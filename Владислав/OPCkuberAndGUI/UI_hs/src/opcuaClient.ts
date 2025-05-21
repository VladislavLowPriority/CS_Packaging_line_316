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

const endpointUrl = 'opc.tcp://WIN-FFD2V15SURQ:4334/UA/OPCUAserverHS' // адрес OPC UA сервера

let client: OPCUAClient
let session: ClientSession | null = null

//Websocket read
export async function readOpcuaStates() {
	if (!session) throw new Error('Session not initialized')

	const tags = {
		start: getNodeIdByTag('start_hs'),
		stop: getNodeIdByTag('stop_hs'),
		reset: getNodeIdByTag('back_to_start'),
	}

	const dataValues = await session.readVariableValue([
		tags.start,
		tags.stop,
		tags.reset,
	])

	return {
		start: dataValues[0].value.value,
		stop: dataValues[1].value.value,
		reset: dataValues[2].value.value,
	}
}
export async function connectToOpcuaServer() {
	client = OPCUAClient.create({
		securityMode: MessageSecurityMode.None,
		securityPolicy: SecurityPolicy.None,
		endpointMustExist: false,
	})

	await client.connect(endpointUrl)
	session = await client.createSession()
	console.log('OPC UA session created')
}

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

function getNodeIdByTag(tag: string): NodeId {
	const map: Record<string, { ns: number; i: number }> = {
		start_hs: { ns: 1, i: 51 },
		stop_hs: { ns: 1, i: 52 },
		back_to_start: { ns: 1, i: 53 },
	}

	const nodeInfo = map[tag]
	if (!nodeInfo) throw new Error(`Unknown tag: ${tag}`)

	return new NodeId(NodeIdType.NUMERIC, nodeInfo.i, nodeInfo.ns)
}
