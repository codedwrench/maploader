package robot

var robots []Robot = []Robot{
	/*
		{model: "your.hostname", // Debug config
		    Serials:          []string{},
			MapFolder:        "folder",
			mapFiles:         []string{},
			restartProcesses: []Process{}},
	*/
	// Serials are read from /oem/midea/device.sn8
	{model: "midea", // Eureka J12 Ultra
		Serials:	[]string{"750Y000D", "750Y000J"},
		// Single folder as opposed to original repo, these vacuums have 3 slots and backup files, in the same folder.
		MapFolder:	"/oem/slam",
		// Let's go for slot 1 here, we'll deduce the other slots from that and also the backup files
		mapFiles:	[]string{"OriginSlamMap01.map.gz",
							 "auto_room_map01.bin.gz",
							 "room_beautify_map01.bin.gz",
							 "room_map01.bin.gz",
							 "carpet_map01.bin.gz",
							 "room_cfg01.settings"
							 "room_info1.yml",
							 "carpet_info1.yml",
							 "bridge_1.settings",
							 "charge_dock_1.settings",
							 "cruise_point_1.settings",
							 "curtain_1.settings",
							 "forb_zone_1.settings",
							 "manual_lines_1.settings",
							 "threshold_1.settings",
							 "histotorical_trajectory1.yml"
							 "room_map_to_planner.bin"} // this one, we just need to delete it
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "midea", // Eureka J15 Ultra
		Serials:	[]string{"750Y000V", "750Y000W"},
		// Single folder as opposed to original repo, these vacuums have 3 slots and backup files, in the same folder.
		MapFolder:	"/oem/slam",
		// Let's go for slot 1 here, we'll deduce the other slots from that and also the backup files
		mapFiles:	[]string{"OriginSlamMap01.map.gz",
							 "auto_room_map01.bin.gz",
							 "room_beautify_map01.bin.gz",
							 "room_map01.bin.gz",
							 "carpet_map01.bin.gz",
							 "room_cfg01.settings"
							 "room_info1.yml",
							 "carpet_info1.yml",
							 "bridge_1.settings",
							 "charge_dock_1.settings",
							 "cruise_point_1.settings",
							 "curtain_1.settings",
							 "forb_zone_1.settings",
							 "manual_lines_1.settings",
							 "threshold_1.settings",
							 "histotorical_trajectory1.yml"
							 "room_map_to_planner.bin"} // this one, we just need to delete it
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "midea", // Eureka J15 Pro Ultra
		Serials:          []string{"750Y000R", "750Y000U", "750Y000Y"},
		// Single folder as opposed to original repo, these vacuums have 3 slots and backup files, in the same folder.
		MapFolder:	"/oem/slam",
		// Let's go for slot 1 here, we'll deduce the other slots from that and also the backup files
		mapFiles:	[]string{"OriginSlamMap01.map.gz",
							 "auto_room_map01.bin.gz",
							 "room_beautify_map01.bin.gz",
							 "room_map01.bin.gz",
							 "carpet_map01.bin.gz",
							 "room_cfg01.settings"
							 "room_info1.yml",
							 "carpet_info1.yml",
							 "bridge_1.settings",
							 "charge_dock_1.settings",
							 "cruise_point_1.settings",
							 "curtain_1.settings",
							 "forb_zone_1.settings",
							 "manual_lines_1.settings",
							 "threshold_1.settings",
							 "histotorical_trajectory1.yml"
							 "room_map_to_planner.bin"} // this one, we just need to delete it
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "midea", // Eureka J15 Max Ultra
		Serials:          []string{"750Y0014", "750Y0013"},
		// Single folder as opposed to original repo, these vacuums have 3 slots and backup files, in the same folder.
		MapFolder:	"/oem/slam",
		// Let's go for slot 1 here, we'll deduce the other slots from that and also the backup files
		mapFiles:	[]string{"OriginSlamMap01.map.gz",
							 "auto_room_map01.bin.gz",
							 "room_beautify_map01.bin.gz",
							 "room_map01.bin.gz",
							 "carpet_map01.bin.gz",
							 "room_cfg01.settings"
							 "room_info1.yml",
							 "carpet_info1.yml",
							 "bridge_1.settings",
							 "charge_dock_1.settings",
							 "cruise_point_1.settings",
							 "curtain_1.settings",
							 "forb_zone_1.settings",
							 "manual_lines_1.settings",
							 "threshold_1.settings",
							 "histotorical_trajectory1.yml"
							 "room_map_to_planner.bin"} // this one, we just need to delete it
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
}
