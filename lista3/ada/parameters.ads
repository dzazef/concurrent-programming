package Parameters is
    Start_Delay             : constant := 2000;
    
    Max_Tasks               : constant := 60;
    Storage_Capacity        : constant := 60;
    
    Ceo_Speed               : constant := 0.1;
    Worker_Speed            : constant := 0.2;
    Client_Speed            : constant := 0.3;
    Machine_Speed           : constant := 0.2;
    Service_Worker_Speed    : constant := 0.3;
    
    Workers                 : constant := 10;
    Clients                 : constant := 10;
    Machines                : constant := 3;
    Service_Workers         : constant := 2;
    
    Max_Arguments           : constant := 1000;
    
    Max_Workers_In_Queue    : constant := 10000;
    
    Break_Probability       : constant := 0.5;

    Worker_Impatient_Delay  : constant := 0.2;
end Parameters;
