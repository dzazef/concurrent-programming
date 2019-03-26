with Ada.Text_IO; use Ada.Text_IO;
with Parameters;

procedure Main is
    type Boolean is (True, False);
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --     
    type Task_Record is record
        Arg1    : Integer;
        Arg2    : Integer;
        Op      : Integer;
    end record;

    type Task_Array_Type is array (0 .. Parameters.Max_Tasks-1) of Task_Record;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    protected type Task_FIFO_Type is
        entry Push (Item : in Task_Record);
        entry Pop (Item: out Task_Record);
        procedure Print_All;
    private 
        Head    : Natural range 0 .. Parameters.Max_Tasks := 0;
        Tail    : Natural range 0 .. Parameters.Max_Tasks := 0;
        Length  : Natural range 0 .. Parameters.Max_Tasks := 0;
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
        begin
            for I in 0 .. Length-1 loop
                Put_Line ("Arg1: " & Integer'Image(Data(Head+I).Arg1) & " Arg2: " & Integer'Image(Data(Head+I).Arg2) & " Op: " & Integer'Image(Data(Head+I).Op));
            end loop;
        end Print_All;
    end Task_FIFO_Type;

    Task_FIFO : Task_FIFO_Type;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    type Product_Record is record
        Id    : Integer;
        Value    : Integer;
    end record;

    type Product_Array_Type is array (0 .. Parameters.Storage_Capacity-1) of Product_Record;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    protected type Product_FIFO_Type is
        entry Push (Item : in Product_Record);
    entry Pop (Item: out Product_Record);
    procedure Print_All;
    private 
        Head    : Natural range 0 .. Parameters.Storage_Capacity := 0;
        Tail    : Natural range 0 .. Parameters.Storage_Capacity := 0;
        Length  : Natural range 0 .. Parameters.Storage_Capacity := 0;
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
        begin
            for I in 0 .. Length-1 loop
                Put_Line ("Product ID: " & Integer'Image(Data(Head+I).Id) & " Value: " & Integer'Image(Data(Head+I).Value));
            end loop;
        end Print_All;
    end Product_FIFO_Type;

    Product_FIFO : Product_FIFO_Type;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    task type CEO_Type;

    task body CEO_Type is 
    begin
        null;
    end CEO_Type;

    CEO : CEO_Type;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    task type Worker_Type is
        entry Do_Task (Line : in String);
    end Worker_Type;

    task body Worker_Type is 
    begin
        accept Do_Task (Line : in String) do
            Put_Line (Line);
        end Do_Task;
    end Worker_Type;

    type Worker_Array_Type is array (0 .. Parameters.Workers-1) of Worker_Type;
    Worker_Array : Worker_Array_Type;
-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --     
    task type Client_Type;

    task body Client_Type is
    begin
        null;
    end Client_Type;

    type Client_Array_Type is array (0 .. Parameters.Clients-1) of Client_Type;


    
begin
    Product_FIFO.Push ((1, 2));
    Product_FIFO.Print_All;
    Product_FIFO.Push ((2, 2));
    Product_FIFO.Print_All;
    Product_FIFO.Push ((3, 2));
    Product_FIFO.Print_All;
    Product_FIFO.Push ((4, 2));
    Product_FIFO.Print_All;
    Product_FIFO.Push ((5, 2));
    Product_FIFO.Print_All;
    Worker_Array(0).Do_Task("End.");
end Main;
