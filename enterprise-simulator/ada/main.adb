with Ada.Text_IO; use Ada.Text_IO;
with Ada.Numerics.discrete_Random;
with Parameters;

procedure Main is
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --     
    type Mode_Type is (Talkative, Silent);
    subtype Operator is Integer range 0 .. 3;
    subtype Task_Int is Integer range 0 .. Parameters.Max_Tasks;
    subtype Product_Int is Integer range 0 .. Parameters.Storage_Capacity;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    Mode : Mode_Type := Silent;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    function Op_To_String(o: in Operator) return String is
    begin
        case o is
            when 0 => return "+";
            when 1 => return "-";
            when 2 => return "*";
            when 3 => return "div";
        end case;
    end Op_To_String;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
-- Random number generator    
    subtype Rand_Gen_Range is Integer range 1 .. Parameters.Max_Arguments;
    package Rand_Int is new Ada.Numerics.Discrete_Random(Rand_Gen_Range);
    Generator : Rand_Int.Generator;

    function Gen_Int (n: in Rand_Gen_Range) return Integer is
    begin
        return Rand_Int.Random(Generator) mod n;
    end Gen_Int;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --     
-- Task record and task array declaration    
    type Task_Record is record
        Arg1    : Integer;
        Arg2    : Integer;
        Op      : Operator;
    end record;

    type Task_Array_Type is array (0 .. Parameters.Max_Tasks-1) of Task_Record;
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
                Put_Line (Integer'Image(Data(Index).Arg1) & " " & Op_To_String(Data(Index).Op) & Integer'Image(Data(Index).Arg2));
            end loop;
        end Print_All;
    end Task_FIFO_Type;

    Task_FIFO : Task_FIFO_Type;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
-- Declaration of product record
    type Product_Record is record
        Id         : Integer;
        Value      : Integer;
    end record;

    type Product_Array_Type is array (0 .. Parameters.Storage_Capacity-1) of Product_Record;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
-- Declaration of protected FIFO of product records    
    protected type Product_FIFO_Type is
        entry Push (Item : in Product_Record);
    entry Pop (Item: out Product_Record);
    procedure Print_All;
    private 
        Head    : Product_Int := 0;
        Tail    : Product_Int := 0;
        Length  : Product_Int := 0;
        Data    : Product_Array_Type;
    end Product_FIFO_Type;

    protected body Product_FIFO_Type is
        entry Push (Item : in Product_Record)
            when Length < Parameters.Storage_Capacity is
            begin
                Data(Tail) := Item;
                Tail := (Tail + 1) mod Parameters.Storage_Capacity;
                Length := Length + 1;
        end Push;
        entry Pop (Item : out Product_Record)
            when Length > 0 is
            begin
                Item := Data(Head);
                Head := (Head + 1) mod Parameters.Storage_Capacity;
                Length := Length - 1;
        end Pop;
        procedure Print_All is
            Index : Product_Int;
        begin                
            for I in 0 .. Length-1 loop
                Index := ((Head + I) mod Parameters.Storage_Capacity);
                Put_Line ("Product ID: " & Integer'Image(Data(Index).Id) & " Value: " & Integer'Image(Data(Index).Value));
            end loop;
        end Print_All;
    end Product_FIFO_Type;

    Product_FIFO : Product_FIFO_Type;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    function Calculate(t: in Task_Record) return Integer is
    begin
        case t.Op is
            when 0 => return t.Arg1 + t.Arg2;
            when 1 => return t.Arg1 - t.Arg2;
            when 2 => return t.Arg1 * t.Arg2;
            when 3 => return t.Arg1 / t.Arg2;
        end case;
    end Calculate;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    task type CEO_Type;

    task body CEO_Type is 
        Arg1, Arg2, Op : Integer;
    begin
        loop
            delay Parameters.Ceo_Speed;
            Arg1 := Gen_Int(Parameters.Max_Arguments);
            Arg2 := Gen_Int(Parameters.Max_Arguments);
            Op := Gen_Int(4);
            Task_FIFO.Push((Arg1, Arg2, Op));
            if Mode=Talkative then 
                Put_Line ("CEO made task:" & Integer'Image(Arg1) & " " & Op_To_String(Op) &  Integer'Image(Arg2));
            end if;
        end loop;
    end CEO_Type;

    CEO : CEO_Type;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    task type Worker_Type;
    
    task body Worker_Type is
        Task_TODO : Task_Record;
        Result : Integer;
        ID : Integer;
    begin
        loop
            delay Parameters.Worker_Speed;
            Task_FIFO.Pop(Task_TODO);
            Result := Calculate(Task_TODO);
            ID := Gen_Int(Parameters.Max_Arguments);
            Product_FIFO.Push((ID, Result));
            if Mode=Talkative then
                Put_Line ("Worker made task no" & Integer'Image(ID) & " from" & Integer'Image(Task_TODO.Arg1) & " " & Op_To_String(Task_TODO.Op) & Integer'Image(Task_TODO.Arg2) & " with result " & Integer'Image(Result));
            end if;
        end loop;
    end Worker_Type;

    type Worker_Array_Type is array (0 .. Parameters.Workers-1) of Worker_Type;
    
    Worker_Array : Worker_Array_Type;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --     
    task type Client_Type;

    task body Client_Type is
        Product : Product_Record;
    begin
        loop
            delay Parameters.Client_Speed;
            Product_FIFO.Pop(Product);
            if Mode=Talkative then
                Put_Line("Client took product no" & Integer'Image(Product.Id) & " from storage.");
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
    end Print_Help;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    Output : String(1 .. 20);
    Length : Natural;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
begin
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
            end if;
        end loop;
    else
        null;
    end if;
end Main;
