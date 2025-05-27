// opcuaLocalServerClient.ts
import {
	AttributeIds,
	ClientSession,
	DataType,
	DataValue,
	NodeId,
	NodeIdType,
	OPCUAClient,
	StatusCodes,
	Variant,
	WriteValueOptions,
} from 'node-opcua'

const LOCAL_ENDPOINT = 'opc.tcp://WIN-FFD2V15SURQ:4334/UA/OPCUAserverHS'

let localClient: OPCUAClient
let localSession: ClientSession | null = null

let isKubeConnected = false

export async function connectToKuberOpcuaServer() {
	try {
		localClient = OPCUAClient.create({
			/*...*/
		})
		await localClient.connect(LOCAL_ENDPOINT)
		localSession = await localClient.createSession()
		isKubeConnected = true
		console.log('Connected to Kuber OPC UA')
	} catch (err) {
		isKubeConnected = false
		console.error('Kuber connection error:', err)
		throw err
	}
}

export function isKuberConnected(): boolean {
	return isKubeConnected
}

export async function writeToLocalServer(tag: string, value: any, i: number) {
	if (!localSession) throw new Error('Local session not initialized')

	const nodeId = new NodeId(NodeIdType.NUMERIC, i, 1) // ns=1, i совпадает

	const dataValue = new DataValue({
		value: new Variant({ dataType: DataType.Boolean, value }), // bool
		statusCode: StatusCodes.Good,
	})

	const writeValue: WriteValueOptions = {
		nodeId: nodeId,
		attributeId: AttributeIds.Value,
		value: dataValue,
	}

	const statusCode = await localSession.write(writeValue)

	if (statusCode !== StatusCodes.Good) {
		console.error(
			`❌ Failed to write ${tag} to local server:`,
			statusCode.toString()
		)
		return
	}
	console.log(`✅ Written to local OPC UA [${tag} -> i=${i}, ns=1]: ${value}`)
}
