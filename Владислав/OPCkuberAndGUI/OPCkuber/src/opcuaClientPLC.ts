import {
	AttributeIds,
	ClientSession,
	DataType,
	DataValue,
	MessageSecurityMode,
	NodeId,
	NodeIdType,
	OPCUAClient,
	SecurityPolicy,
	StatusCodes,
	Variant,
	WriteValueOptions,
} from 'node-opcua'

const endpointUrl = 'opc.tcp://10.160.160.61:4840'

let client: OPCUAClient
let session: ClientSession | null = null

// Полная карта всех тегов сервера
const TAG_MAP: Record<string, { ns: number; i: number }> = {
	// Основные команды

	// Processing Station Inputs
	processing_input_4_workpiece_detected: { ns: 4, i: 5 },
	processing_input_2_workpiece_silver: { ns: 4, i: 7 },
	processing_input_5_carousel_init: { ns: 4, i: 3 },
	processing_input_6_hole_detected: { ns: 4, i: 4 },
	processing_input_7_workpiece_not_black: { ns: 4, i: 6 },

	// Processing Station Outputs
	processing_output_0_drill: { ns: 4, i: 12 },
	processing_output_1_rotate_carousel: { ns: 4, i: 13 },
	processing_output_2_drill_down: { ns: 4, i: 14 },
	processing_output_3_drill_up: { ns: 4, i: 15 },
	processing_output_4_fix_workpiece: { ns: 4, i: 16 },
	processing_output_5_detect_hole: { ns: 4, i: 17 },

	// Handling and Packing Inputs
	handling_input_0_workpiece_pushed: { ns: 4, i: 29 },
	handling_input_1_grippe_at_right: { ns: 4, i: 32 },
	handling_input_2_gripper_at_start: { ns: 4, i: 31 },
	handling_input_3_gripper_down_pack_lvl: { ns: 4, i: 33 },
	packing_input_7_pack_turned_on: { ns: 4, i: 42 },

	// Handling and Packing Outputs
	handling_output_0_to_green: { ns: 4, i: 34 },
	handling_output_1_to_yellow: { ns: 4, i: 35 },
	handling_output_2_to_red: { ns: 4, i: 36 },
	handling_output_3_gripper_to_right: { ns: 4, i: 37 },
	handling_output_4_gripper_to_left: { ns: 4, i: 38 },
	handling_output_5_gripper_to_down: { ns: 4, i: 39 },
	handling_output_6_gripper_to_open: { ns: 4, i: 40 },
	handling_output_7_gripper_push_workpiece: { ns: 4, i: 41 },
	packing_output_4_push_box: { ns: 4, i: 43 },
	packing_output_5_fix_box_upper_side: { ns: 4, i: 44 },
	packing_output_6_fix_box_tongue: { ns: 4, i: 45 },
	packing_output_7_pack_box: { ns: 4, i: 46 },

	// Sorting Station Inputs
	sorting_input_3_box_on_conveyor: { ns: 4, i: 9 },
	sorting_input_4_box_is_down: { ns: 4, i: 10 },

	// Sorting Station Outputs
	sorting_output_0_move_conveyor_right: { ns: 4, i: 19 },
	sorting_output_1_move_conveyor_left: { ns: 4, i: 20 },
	sorting_output_2_push_silver_workpiece: { ns: 4, i: 21 },
	sorting_output_3_push_red_workpiece: { ns: 4, i: 22 },
}
const tagsToRead = [
	'processing_output_1_rotate_carousel',
	'sorting_output_0_move_conveyor_right',
	'sorting_output_3_push_red_workpiece',
] //так можем читать нужные теги

export async function readOpcuaStates(): Promise<Record<string, any>> {
	if (!session) throw new Error('Session not initialized')

	const nodeIds = tagsToRead.map(tag => getNodeIdByTag(tag))
	const startTime = Date.now()

	try {
		const dataValues = await session.readVariableValue(nodeIds)

		if (!dataValues || dataValues.length !== nodeIds.length) {
			throw new Error(
				'Некорректное количество возвращённых значений от OPC UA сервера'
			)
		}

		const tagValues: Record<string, any> = {}
		tagsToRead.forEach((tag, index) => {
			const value = dataValues[index]?.value?.value
			tagValues[tag] = value
		})

		const executionTime = Date.now() - startTime

		// Вывод в консоль: время + значения тегов
		const timeString = new Date().toLocaleTimeString()
		const tagStatus = tagsToRead
			.map(tag => `${tag}: ${tagValues[tag]}`)
			.join(', ')
		console.log(` `, tagStatus)

		return tagValues
	} catch (error) {
		console.error('Error reading OPCUA tags:', error)
		throw error
	}
}

// Подключение к серверу (остается без изменений)
let isConnected = false

export async function connectToOpcuaServer() {
	try {
		client = OPCUAClient.create({
			securityMode: MessageSecurityMode.None,
			securityPolicy: SecurityPolicy.None,
			endpointMustExist: false,
		})

		await client.connect(endpointUrl)
		session = await client.createSession()
		isConnected = true
		console.log('OPC UA session created')
	} catch (error) {
		isConnected = false
		console.error('Failed to connect to OPC UA server:', error)
		throw error
	}
}

export function isOpcuaConnected(): boolean {
	return isConnected
}

// Вспомогательная функция для получения NodeId
function getNodeIdByTag(tag: string): NodeId {
	const nodeInfo = TAG_MAP[tag]
	if (!nodeInfo) throw new Error(`Unknown tag: ${tag}`)
	return new NodeId(NodeIdType.NUMERIC, nodeInfo.i, nodeInfo.ns)
}

// Функция отправки команд (остается без изменений)
export async function sendCommandToOpcua(tag: string, value: boolean) {
	if (!session) throw new Error('Session not initialized')

	const nodeId = getNodeIdByTag(tag)

	const dataValue = new DataValue({
		value: new Variant({ dataType: DataType.Boolean, value }),
		statusCode: StatusCodes.Good,
	})

	const writeValue: WriteValueOptions = {
		nodeId: nodeId,
		attributeId: AttributeIds.Value,
		value: dataValue,
	}

	const statusCode = await session.write(writeValue)

	if (statusCode !== StatusCodes.Good) {
		throw new Error(`Failed to write value: ${statusCode.toString()}`)
	}

	console.log(
		`🛰 Sent to OPC UA: ${tag} = ${value}, status: ${statusCode.toString()}`
	)
}
