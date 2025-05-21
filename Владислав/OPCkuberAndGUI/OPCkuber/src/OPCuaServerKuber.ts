const { OPCUAServer, Variant, DataType, StatusCodes } = require('node-opcua')

export async function startOpcuaServer() {
	const server = new OPCUAServer({
		port: 4334,
		resourcePath: '/UA/OPCUAserverHS',
		buildInfo: {
			productName: 'OPCuaServerHS',
			buildNumber: '1',
			buildDate: new Date(),
		},
	})

	await server.initialize()

	const addressSpace = server.engine.addressSpace
	const namespace = addressSpace.getOwnNamespace()

	// Объявление переменных
	let Start_hs = false
	let stop_hs = false
	let back_to_start = false

	// Processing Station
	let processing_input_1_workpiece_detected = false
	let processing_input_2_workpiece_silver = false
	let processing_input_5_carousel_init = false
	let processing_input_6_hole_detected = false
	let processing_input_7_workpiece_not_black = false
	let processing_output_0_drill = false
	let processing_output_1_rotate_carousel = false
	let processing_output_2_drill_down = false
	let processing_output_3_drill_up = false
	let processing_output_4_fix_workpiece = false
	let processing_output_5_detect_hole = false

	// Handling and Packing Station
	let handling_input_0_workpiece_pushed = false
	let handling_input_1_grippe_at_right = false
	let handling_input_2_gripper_at_start = false
	let handling_input_3_gripper_down_pack_lvl = false
	let packing_input_7_pack_turned_on = false
	let handling_output_0_to_green = false
	let handling_output_1_to_yellow = false
	let handling_output_2_to_red = false
	let handling_output_3_gripper_to_right = false
	let handling_output_4_gripper_to_left = false
	let handling_output_5_gripper_to_down = false
	let handling_output_6_gripper_to_open = false
	let handling_output_7_gripper_push_workpiece = false
	let packing_output_4_push_box = false
	let packing_output_5_fix_box_upper_side = false
	let packing_output_6_fix_box_tongue = false
	let packing_output_7_pack_box = false

	// Sorting Station
	let sorting_input_3_box_on_conveyor = false
	let sorting_input_4_box_is_down = false
	let sorting_output_0_move_conveyor_right = false
	let sorting_output_1_move_conveyor_left = false
	let sorting_output_2_push_silver_workpiece = false
	let sorting_output_3_push_red_workpiece = false

	// Создание основного устройства
	const device = namespace.addObject({
		organizedBy: addressSpace.rootFolder.objects,
		browseName: 'PLC_Controller',
	})

	// Основные команды start/stop/reset
	namespace.addVariable({
		componentOf: device,
		browseName: 'start_hs',
		nodeId: 'ns=1;i=51',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () => new Variant({ dataType: DataType.Boolean, value: Start_hs }),
			set: (variant: { value: boolean }) => {
				Start_hs = variant.value
				console.log('start_hs:', Start_hs)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: device,
		browseName: 'stop_hs',
		nodeId: 'ns=1;i=52',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () => new Variant({ dataType: DataType.Boolean, value: stop_hs }),
			set: (variant: { value: boolean }) => {
				stop_hs = variant.value
				console.log('stop_hs:', stop_hs)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: device,
		browseName: 'back_to_start',
		nodeId: 'ns=1;i=53',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({ dataType: DataType.Boolean, value: back_to_start }),
			set: (variant: { value: boolean }) => {
				back_to_start = variant.value
				console.log('back_to_start:', back_to_start)
				return StatusCodes.Good
			},
		},
	})

	// Processing Station
	const processingStation = namespace.addObject({
		componentOf: device,
		browseName: 'ProcessingStation',
	})

	// Inputs
	const processingInputs = namespace.addObject({
		componentOf: processingStation,
		browseName: 'Inputs',
	})

	namespace.addVariable({
		componentOf: processingInputs,
		browseName: 'processing_input_1_workpiece_detected',
		nodeId: 'ns=1;i=5',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: processing_input_1_workpiece_detected,
				}),
		},
	})

	namespace.addVariable({
		componentOf: processingInputs,
		browseName: 'processing_input_2_workpiece_silver',
		nodeId: 'ns=1;i=7',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: processing_input_2_workpiece_silver,
				}),
		},
	})

	namespace.addVariable({
		componentOf: processingInputs,
		browseName: 'processing_input_5_carousel_init',
		nodeId: 'ns=1;i=3',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: processing_input_5_carousel_init,
				}),
		},
	})

	namespace.addVariable({
		componentOf: processingInputs,
		browseName: 'processing_input_6_hole_detected',
		nodeId: 'ns=1;i=4',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: processing_input_6_hole_detected,
				}),
		},
	})

	namespace.addVariable({
		componentOf: processingInputs,
		browseName: 'processing_input_7_workpiece_not_black',
		nodeId: 'ns=1;i=6',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: processing_input_7_workpiece_not_black,
				}),
		},
	})

	// Outputs
	const processingOutputs = namespace.addObject({
		componentOf: processingStation,
		browseName: 'Outputs',
	})

	namespace.addVariable({
		componentOf: processingOutputs,
		browseName: 'processing_output_0_drill',
		nodeId: 'ns=1;i=12',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: processing_output_0_drill,
				}),
			set: (variant: { value: boolean }) => {
				processing_output_0_drill = variant.value
				console.log('processing_output_0_drill:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: processingOutputs,
		browseName: 'processing_output_1_rotate_carousel',
		nodeId: 'ns=1;i=13',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: processing_output_1_rotate_carousel,
				}),
			set: (variant: { value: boolean }) => {
				processing_output_1_rotate_carousel = variant.value
				console.log('processing_output_1_rotate_carousel:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: processingOutputs,
		browseName: 'processing_output_2_drill_down',
		nodeId: 'ns=1;i=14',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: processing_output_2_drill_down,
				}),
			set: (variant: { value: boolean }) => {
				processing_output_2_drill_down = variant.value
				console.log('processing_output_2_drill_down:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: processingOutputs,
		browseName: 'processing_output_3_drill_up',
		nodeId: 'ns=1;i=15',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: processing_output_3_drill_up,
				}),
			set: (variant: { value: boolean }) => {
				processing_output_3_drill_up = variant.value
				console.log('processing_output_3_drill_up:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: processingOutputs,
		browseName: 'processing_output_4_fix_workpiece',
		nodeId: 'ns=1;i=16',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: processing_output_4_fix_workpiece,
				}),
			set: (variant: { value: boolean }) => {
				processing_output_4_fix_workpiece = variant.value
				console.log('processing_output_4_fix_workpiece:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: processingOutputs,
		browseName: 'processing_output_5_detect_hole',
		nodeId: 'ns=1;i=17',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: processing_output_5_detect_hole,
				}),
			set: (variant: { value: boolean }) => {
				processing_output_5_detect_hole = variant.value
				console.log('processing_output_5_detect_hole:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	// Handling and Packing Station
	const handlingPacking = namespace.addObject({
		componentOf: device,
		browseName: 'HandlingPacking',
	})

	// Inputs
	const handlingInputs = namespace.addObject({
		componentOf: handlingPacking,
		browseName: 'Inputs',
	})

	namespace.addVariable({
		componentOf: handlingInputs,
		browseName: 'handling_input_0_workpiece_pushed',
		nodeId: 'ns=1;i=29',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: handling_input_0_workpiece_pushed,
				}),
		},
	})

	namespace.addVariable({
		componentOf: handlingInputs,
		browseName: 'handling_input_1_grippe_at_right',
		nodeId: 'ns=1;i=32',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: handling_input_1_grippe_at_right,
				}),
		},
	})

	namespace.addVariable({
		componentOf: handlingInputs,
		browseName: 'handling_input_2_gripper_at_start',
		nodeId: 'ns=1;i=31',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: handling_input_2_gripper_at_start,
				}),
		},
	})

	namespace.addVariable({
		componentOf: handlingInputs,
		browseName: 'handling_input_3_gripper_down_pack_lvl',
		nodeId: 'ns=1;i=33',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: handling_input_3_gripper_down_pack_lvl,
				}),
		},
	})

	namespace.addVariable({
		componentOf: handlingInputs,
		browseName: 'packing_input_7_pack_turned_on',
		nodeId: 'ns=1;i=42',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: packing_input_7_pack_turned_on,
				}),
		},
	})

	// Outputs
	const handlingOutputs = namespace.addObject({
		componentOf: handlingPacking,
		browseName: 'Outputs',
	})

	namespace.addVariable({
		componentOf: handlingOutputs,
		browseName: 'handling_output_0_to_green',
		nodeId: 'ns=1;i=34',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: handling_output_0_to_green,
				}),
			set: (variant: { value: boolean }) => {
				handling_output_0_to_green = variant.value
				console.log('handling_output_0_to_green:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: handlingOutputs,
		browseName: 'handling_output_1_to_yellow',
		nodeId: 'ns=1;i=35',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: handling_output_1_to_yellow,
				}),
			set: (variant: { value: boolean }) => {
				handling_output_1_to_yellow = variant.value
				console.log('handling_output_1_to_yellow:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: handlingOutputs,
		browseName: 'handling_output_2_to_red',
		nodeId: 'ns=1;i=36',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: handling_output_2_to_red,
				}),
			set: (variant: { value: boolean }) => {
				handling_output_2_to_red = variant.value
				console.log('handling_output_2_to_red:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: handlingOutputs,
		browseName: 'handling_output_3_gripper_to_right',
		nodeId: 'ns=1;i=37',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: handling_output_3_gripper_to_right,
				}),
			set: (variant: { value: boolean }) => {
				handling_output_3_gripper_to_right = variant.value
				console.log('handling_output_3_gripper_to_right:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: handlingOutputs,
		browseName: 'handling_output_4_gripper_to_left',
		nodeId: 'ns=1;i=38',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: handling_output_4_gripper_to_left,
				}),
			set: (variant: { value: boolean }) => {
				handling_output_4_gripper_to_left = variant.value
				console.log('handling_output_4_gripper_to_left:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: handlingOutputs,
		browseName: 'handling_output_5_gripper_to_down',
		nodeId: 'ns=1;i=39',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: handling_output_5_gripper_to_down,
				}),
			set: (variant: { value: boolean }) => {
				handling_output_5_gripper_to_down = variant.value
				console.log('handling_output_5_gripper_to_down:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: handlingOutputs,
		browseName: 'handling_output_6_gripper_to_open',
		nodeId: 'ns=1;i=40',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: handling_output_6_gripper_to_open,
				}),
			set: (variant: { value: boolean }) => {
				handling_output_6_gripper_to_open = variant.value
				console.log('handling_output_6_gripper_to_open:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: handlingOutputs,
		browseName: 'handling_output_7_gripper_push_workpiece',
		nodeId: 'ns=1;i=41',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: handling_output_7_gripper_push_workpiece,
				}),
			set: (variant: { value: boolean }) => {
				handling_output_7_gripper_push_workpiece = variant.value
				console.log('handling_output_7_gripper_push_workpiece:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: handlingOutputs,
		browseName: 'packing_output_4_push_box',
		nodeId: 'ns=1;i=43',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: packing_output_4_push_box,
				}),
			set: (variant: { value: boolean }) => {
				packing_output_4_push_box = variant.value
				console.log('packing_output_4_push_box:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: handlingOutputs,
		browseName: 'packing_output_5_fix_box_upper_side',
		nodeId: 'ns=1;i=44',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: packing_output_5_fix_box_upper_side,
				}),
			set: (variant: { value: boolean }) => {
				packing_output_5_fix_box_upper_side = variant.value
				console.log('packing_output_5_fix_box_upper_side:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: handlingOutputs,
		browseName: 'packing_output_6_fix_box_tongue',
		nodeId: 'ns=1;i=45',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: packing_output_6_fix_box_tongue,
				}),
			set: (variant: { value: boolean }) => {
				packing_output_6_fix_box_tongue = variant.value
				console.log('packing_output_6_fix_box_tongue:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: handlingOutputs,
		browseName: 'packing_output_7_pack_box',
		nodeId: 'ns=1;i=46',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: packing_output_7_pack_box,
				}),
			set: (variant: { value: boolean }) => {
				packing_output_7_pack_box = variant.value
				console.log('packing_output_7_pack_box:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	// Sorting Station
	const sortingStation = namespace.addObject({
		componentOf: device,
		browseName: 'SortingStation',
	})

	// Inputs
	const sortingInputs = namespace.addObject({
		componentOf: sortingStation,
		browseName: 'Inputs',
	})

	namespace.addVariable({
		componentOf: sortingInputs,
		browseName: 'sorting_input_3_box_on_conveyor',
		nodeId: 'ns=1;i=9',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: sorting_input_3_box_on_conveyor,
				}),
		},
	})

	namespace.addVariable({
		componentOf: sortingInputs,
		browseName: 'sorting_input_4_box_is_down',
		nodeId: 'ns=1;i=10',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: sorting_input_4_box_is_down,
				}),
		},
	})

	// Outputs
	const sortingOutputs = namespace.addObject({
		componentOf: sortingStation,
		browseName: 'Outputs',
	})

	namespace.addVariable({
		componentOf: sortingOutputs,
		browseName: 'sorting_output_0_move_conveyor_right',
		nodeId: 'ns=1;i=19',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: sorting_output_0_move_conveyor_right,
				}),
			set: (variant: { value: boolean }) => {
				sorting_output_0_move_conveyor_right = variant.value
				console.log('sorting_output_0_move_conveyor_right:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: sortingOutputs,
		browseName: 'sorting_output_1_move_conveyor_left',
		nodeId: 'ns=1;i=20',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: sorting_output_1_move_conveyor_left,
				}),
			set: (variant: { value: boolean }) => {
				sorting_output_1_move_conveyor_left = variant.value
				console.log('sorting_output_1_move_conveyor_left:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: sortingOutputs,
		browseName: 'sorting_output_2_push_silver_workpiece',
		nodeId: 'ns=1;i=21',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: sorting_output_2_push_silver_workpiece,
				}),
			set: (variant: { value: boolean }) => {
				sorting_output_2_push_silver_workpiece = variant.value
				console.log('sorting_output_2_push_silver_workpiece:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	namespace.addVariable({
		componentOf: sortingOutputs,
		browseName: 'sorting_output_3_push_red_workpiece',
		nodeId: 'ns=1;i=22',
		dataType: 'Boolean',
		minimumSamplingInterval: 100,
		value: {
			get: () =>
				new Variant({
					dataType: DataType.Boolean,
					value: sorting_output_3_push_red_workpiece,
				}),
			set: (variant: { value: boolean }) => {
				sorting_output_3_push_red_workpiece = variant.value
				console.log('sorting_output_3_push_red_workpiece:', variant.value)
				return StatusCodes.Good
			},
		},
	})

	await server.start()

	console.log('Сервер запущен на порту', server.endpoints[0].port)
	console.log(
		'Endpoint URL:',
		server.endpoints[0].endpointDescriptions()[0].endpointUrl
	)
}
