// main.ts
import {
	connectToKuberOpcuaServer,
	isKuberConnected,
	writeToLocalServer,
} from './opcuaClientKube'
import {
	connectToOpcuaServer,
	isOpcuaConnected,
	readOpcuaStates,
} from './opcuaClientPLC'
import { startOpcuaServer } from './OPCuaServerKuber'

const SYNC_INTERVAL = 2000 // Интервал синхронизации в ms
const RETRY_DELAY = 5000 // Задержка переподключения

// Таблица соответствия тегов PLC -> Kuber
const TAG_MAPPING: Record<string, number> = {
	processing_output_1_rotate_carousel: 13,
	sorting_output_0_move_conveyor_right: 19,
	sorting_output_3_push_red_workpiece: 22,
}

async function syncData() {
	try {
		// 1. Чтение данных из PLC
		const startTime = Date.now()
		const plcData = await readOpcuaStates()

		// 2. Параллельная запись всех значений в Kuber
		const writePromises = Object.entries(plcData).map(async ([tag, value]) => {
			const kuberNodeId = TAG_MAPPING[tag]
			if (kuberNodeId === undefined) return

			try {
				await writeToLocalServer(tag, value, kuberNodeId)
			} catch (err) {
				console.error(`Error writing ${tag}:`, err)
			}
		})

		await Promise.all(writePromises)

		// 3. Логирование производительности
		const syncTime = Date.now() - startTime
		console.log(`Data synced in ${syncTime}ms`)
	} catch (err) {
		console.error('Sync error:', err)
	}
}

async function initializeConnections() {
	try {
		// Запуск сервера Kuber
		await startOpcuaServer()
		console.log('Kuber OPC UA server started')

		// Подключение к PLC
		await connectToOpcuaServer()
		console.log('Connected to PLC OPC UA')

		// Подключение к Kuber
		await connectToKuberOpcuaServer()
		console.log('Connected to Kuber OPC UA')

		// Запуск цикла синхронизации
		setInterval(async () => {
			if (isOpcuaConnected() && isKuberConnected()) {
				await syncData()
			} else {
				console.log('Waiting for connections...')
			}
		}, SYNC_INTERVAL)
	} catch (err) {
		console.error('Initialization failed:', err)
		console.log(`Retrying in ${RETRY_DELAY / 1000} seconds...`)
		setTimeout(initializeConnections, RETRY_DELAY)
	}
}

// Запуск приложения
initializeConnections()
