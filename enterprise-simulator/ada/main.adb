with Ada.Text_IO; use Ada.Text_IO;
with Ada.Numerics.discrete_Random;
with Parameters;

procedure Main is
    type Simulator_Mode is (Talkative, Silent);
    type Worker_Mode_Type is (Patient, Impatient);
    subtype Operator is Integer range 0..1;
    subtype Task_Int is Integer range 0 .. Parameters.Max_Tasks;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
     Mode : Simulator_Mode := Silent;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    procedure Print_Info_If_Talkative(S : in String) is
    begin
        if (Mode = Talkative) then Put_Line(S); end if;
    end Print_Info_If_Talkative;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    function Op_To_String(o: in Operator) return String is
    begin
        case o is
            when 0 => return "+";
            when 1 => return "*";
        end case;
    end Op_To_String;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
-- Random number generator
    subtype Rand_Gen_Range is Integer range 0 .. Parameters.Max_Arguments;
    package Rand_Int is new Ada.Numerics.Discrete_Random(Rand_Gen_Range);
    Generator : Rand_Int.Generator;

    function Gen_Int (n: in Rand_Gen_Range) return Integer is
    begin
        return Rand_Int.Random(Generator) mod n;
    end Gen_Int;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
-- Random worker mode generator
    function Gen_Mode return Worker_Mode_Type is
    begin
        if Gen_Int(2)=1 then
            return Patient;
        else
            return Impatient;
        end if;
    end Gen_Mode;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
-- Task record and task array declaration
    type Task_Record is record
        Arg1    : Integer;
        Arg2    : Integer;
        Op      : Operator;
        Value   : Integer;
    end record;

    type Task_Array_Type is array (0 .. Parameters.Max_Tasks-1) of Task_Record;
    -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    subtype Worker_Num_Type is Integer range 0..Parameters.Workers-1;
    type Worker_Statistic_Array_Type is array (Worker_Num_Type) of Natural;

    protected type Worker_Statistic_Type is
        procedure Increase (ID : in Worker_Num_Type);
        procedure Print_All;
    private
        Data     : Worker_Statistic_Array_Type;
    end Worker_Statistic_Type;

    protected body Worker_Statistic_Type is
        procedure Increase (ID : in Worker_Num_Type) is
        begin
            Data(ID) := Data(ID) + 1;
        end;

        procedure Print_All is
        begin
            for I in Worker_Num_Type'Range loop
                Put_Line("INFO: Worker" & Integer'Image(I) & " done" & Integer'Image(Data(I)) & "tasks.");
            end loop;
        end Print_All;
    end Worker_Statistic_Type;

    Worker_Statistic : Worker_Statistic_Type;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    procedure Calculate_Task(t: in out Task_Record) is
    begin
        case t.Op is
            when 0 => t.Value := t.Arg1 + t.Arg2;
            when 1 => t.Value := t.Arg1 * t.Arg2;
        end case;
    end Calculate_Task;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
-- Declaration of protected FIFO of task records
    protected type Task_FIFO_Type is
        entry Push (Item : in Task_Record);
        entry Pop (Item: out Task_Record);
        procedure Print_All;
    private
        Head    : Task_Int := 0;
        Tail    : Task_Int := 0;
        Length  : Task_Int := 0;
        Data    : Task_Array_Type;
    end Task_FIFO_Type;

    protected body Task_FIFO_Type is
        entry Push (Item : in Task_Record)
            when Length < Parameters.Max_Tasks is
            begin
                Data(Tail) := Item;
                Tail := (Tail + 1) mod Parameters.Max_Tasks;
                Length := Length + 1;
        end Push;
        entry Pop (Item : out Task_Record)
            when Length > 0 is
            begin
                Item := Data(Head);
                Head := (Head + 1) mod Parameters.Max_Tasks;
                Length := Length - 1;
        end Pop;
        procedure Print_All is
            Index : Task_Int;
        begin
            for I in 0 .. Length-1 loop
                Index := ((Head + I) mod Parameters.Max_Tasks);
                Put_Line (Integer'Image(Data(Index).Arg1) & " " & Op_To_String(Data(Index).Op) & Integer'Image(Data(Index).Arg2) & " = " & Integer'Image(Data(Index).Value));
            end loop;
        end Print_All;
    end Task_FIFO_Type;

    Task_FIFO : Task_FIFO_Type;
    Product_FIFO : Task_FIFO_Type;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
-- CEO
    task type CEO_Type is
        entry Start;
    end CEO_Type;

    task body CEO_Type is
        Arg1, Arg2, Op : Integer;
    begin
        accept Start;
        loop
            delay Parameters.Ceo_Speed;
            Arg1 := Gen_Int(Parameters.Max_Arguments);
            Arg2 := Gen_Int(Parameters.Max_Arguments);
            Op := Gen_Int(2);
            Task_FIFO.Push((Arg1, Arg2, Op, 0));
            if Mode=Talkative then
                Put_Line ("CEO made task:" & Integer'Image(Arg1) & " " & Op_To_String(Op) &  Integer'Image(Arg2));
            end if;
        end loop;
    end CEO_Type;

    CEO : CEO_Type;

-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
-- MACHINES

   task type Machine_Type is
        entry Enable(Op1 : in Operator);
        entry Calculate(Item : in out Task_Record);
   end Machine_Type;

    task body Machine_Type is
        Op : Operator;
    begin
        accept Enable(Op1 : in Operator) do
            Op := Op1;
        end Enable;
        loop
            accept Calculate ( Item : in out Task_Record) do
                    Print_Info_If_Talkative("MACHINE(" & Op_To_String(Op) & "): Working on task:"
                             & Integer'Image(Item.Arg1) & " " & Op_To_String(Item.Op) & Integer'Image(Item.Arg2));
                delay Parameters.Machine_Speed;
                Calculate_Task(Item);
            end Calculate;
        end loop;
    end Machine_Type;

    type Machine_Array_Type is array (0 .. Parameters.Machines-1) of Machine_Type;
    type Machine_Set_Type is array (Operator) of Machine_Array_Type;

    Machine_Set : Machine_Set_Type;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    task type Worker_Type is
        entry Start(ID1 : in Natural);
    end Worker_Type;

    task body Worker_Type is
        ID : Natural;
        Worker_Mode : Worker_Mode_Type;
        Current_Task : Task_Record;
    begin
        accept Start(ID1 : in Natural) do
            ID := ID1;
            Worker_Mode := Gen_Mode;
        end Start;
        loop
            delay Parameters.Worker_Speed;
            Task_FIFO.Pop(Current_Task);
            if (Worker_Mode = Patient) then
                    Print_Info_If_Talkative("WORKER" & Natural'Image(ID) & " waiting for machine for task"
                             & Integer'Image(Current_Task.Arg1) & " " & Op_To_String(Current_Task.Op) & Integer'Image(Current_Task.Arg2));
                Machine_Set(Current_Task.Op)(Gen_Int(Parameters.Machines)).Calculate(Current_Task);
            else
                Print_Info_If_Talkative("WORKER" & Natural'Image(ID) & " searching for machine for task"
                             & Integer'Image(Current_Task.Arg1) & " " & Op_To_String(Current_Task.Op) & Integer'Image(Current_Task.Arg2));
                loop
                    select
                        Machine_Set(Current_Task.Op)(Gen_Int(Parameters.Machines)).Calculate(Current_Task);
                        exit;
                    or
                        delay Parameters.Worker_Impatient_Delay;
                    end select;
                end loop;

            end if;
            Product_FIFO.Push(Current_Task);
            Print_Info_If_Talkative("WORKER: Put task" & Integer'Image(Current_Task.Arg1)
                         & " " & Op_To_String(Current_Task.Op) & Integer'Image(Current_Task.Arg2)
                                    & " =" & Integer'Image(Current_Task.Value) &" in storage.");
            Worker_Statistic.Increase(ID);
        end loop;
    end Worker_Type;

    type Worker_Array_Type is array (0 .. Parameters.Workers-1) of Worker_Type;

    Worker_Array : Worker_Array_Type;
 -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    task type Client_Type is
        entry Start;
    end Client_Type;

    task body Client_Type is
        Product : Task_Record;
    begin
        accept Start;
        loop
            delay Parameters.Client_Speed;
            Product_FIFO.Pop(Product);
            if Mode=Talkative then
                Put_Line("CLIENT took product" & Integer'Image(Product.Arg1)
                         & " " & Op_To_String(Product.Op) & Integer'Image(Product.Arg2) & " from storage.");
            end if;
        end loop;
    end Client_Type;

    type Client_Array_Type is array (0 .. Parameters.Clients-1) of Client_Type;

    Client_Array : Client_Array_Type;
 -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    procedure Print_Help is
    begin
        Put_Line ("----HELP----");
        Put_Line ("help - view help");
        Put_Line ("talk - enable talkative mode");
        Put_Line ("tasklist - print active tasks");
        Put_Line ("storage - print storage");
        Put_Line ("statistics - print workers statistics");
    end Print_Help;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    Output : String(1 .. 20);
    Length : Natural;
    Current_Worker_Index : Natural := 0;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
begin
    Ceo.Start;
    for O in Operator'Range loop
        for I in Machine_Set(O)'Range loop
            Machine_Set(O)(I).Enable(O);
        end loop;
    end loop;
    for I in Worker_Array'Range loop
        Worker_Array(I).Start(Current_Worker_Index);
        Current_Worker_Index := Current_Worker_Index + 1;
    end loop;
    for I in Client_Array'Range loop
        Client_Array(I).Start;
    end loop;
    if Mode=Silent then
        Print_Help;
        loop
            Get_Line (Output, Length);
            if Output(1 .. Length)="help" then
                Print_Help;
            elsif Output(1 .. Length)="talk" then
                Mode := Talkative;
            elsif Output(1 .. Length)="tasklist" then
                Task_FIFO.Print_All;
            elsif Output(1 .. Length)="storage" then
                Product_FIFO.Print_All;
            elsif Output(1 .. Length)="statistics" then
                Worker_Statistic.Print_All;
            end if;
        end loop;
   else
        null;
    end if;
end Main;
