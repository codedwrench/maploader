package robot

var robots []Robot = []Robot{
	/*
		{model: "your.hostname", // Debug config
		    Serials:          []string{},
			MapFolders:       []string{},
			mapFiles:         []string{},
			restartProcesses: []Process{}},
	*/
    // Serials are read from /oem/midea/device.sn8
	{model: "midea", // Eureka J12 Ultra
		Serials:          []string{"750Y000D", "750Y000J"}
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "midea", // Eureka J15 Ultra
		Serials:          []string{"750Y000V", "750Y000W"}
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "midea", // Eureka J15 Pro Ultra
	    Serials:          []string{"750Y000R", "750Y000U", "750Y000Y"}
        MapFolders:       []string{"/data/ri", "/data/map", "/data/DivideMap"},
        mapFiles:         []string{"/data/config/ava/mult_map.json"},
        restartProcesses: []Process{MiioClientProcess, AvaProcess}},
    {model: "midea", // Eureka J15 Max Ultra
		Serials:          []string{"750Y0014", "750Y0013"}
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},

}
