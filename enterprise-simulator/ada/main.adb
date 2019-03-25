with Ada.Text_IO; use Ada.Text_IO;
with Parameters;

procedure Main is
    type New_Task is record
        Arg1    : Integer;
        Arg2    : Integer;
        Op      : Integer;
    end record;

    task type Worker is
        entry Do_Task (Line : in String);
    end Worker;

    task body Worker is 
    begin
        accept Do_Task (Line : in String) do
            Put_Line (Line);
        end Do_Task;
    end Worker;

    Worker1 : Worker;
begin
    Worker1.Do_Task("xd");
end Main;
