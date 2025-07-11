# Упаковочная линия 316
## Структура проекта
```
 📦     
 ├── 📂 Documentation               // Документация  
 ├── 📂 Projects                    // Проекты  
 │   ├── 📂 CS                      // Система управления (Golang)  
 │   ├── 📂 EnIPExplorer            // Программа для поиска устройств Ethernet/IP в сети  
 │   ├── 📂 OPCKuber+GUI  
 │   │   ├── 📂 OPCkuber            // Программа для запуска OPC UA сервера в Kubernetes  
 │   │   └── 📂 UI_hs               // Программа реализации минимального пользовательского интерфейса для взаимодействия с СУ  
 │   └── 📂 TiaPortal_EIP_Project   // Проект Tia Portal v16 Ethernet/IP для двух контроллеров S7-1200   
 └── 📜 README.md
```
# Документация проекта
## Авторы
Качалов Владислав  
Полетаев Александр  
Соловьев Михаил  
Васенин Константин  
Чжан Цзяпэн  
Ван Цзыи  
Наставник: **Гайфулин Тимур**  
Руководитель: **Куклин Егор Вадимович**  

## Описание поставленной задачи
Обеспечить сокращение задержек передачи управляющих сигналов между установкой упаковочной станции (СУ 316) и системы управления, развернутой на сервере, путем поиска и интегрирования протокола реального времени (RT). Также необходимо определить подходящий язык программирования и адаптировать под него логику СУ, подготовить интерфейс для реализации функционала старта/стопа и развернуть СУ в Kubernetes.

## Используемое ПО
 - Siemens TIA Portal V16 / SIMATIC STEP 7 incl. Safety and WinCC V16    
 - UaExpert  
 - VisualCode  
 - EnIPExploler  
 - Docker Desktop  
 - Kubernetes  


## Ссылка на репозиторий (2024 год)
[HS Line 316](https://github.com/anvax/HS_line316?tab=readme-ov-file#%D0%B0%D0%B2%D1%82%D0%BE%D1%80%D1%8B)

## Архитектура

![изображение](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Documentation/image.png)

## Ход работы:
## Протокол EtherNet/IP
Была написана тестовая программа СУ на языке Golang, используюящая протокол связи OPCua.  
В ходе тестирования мы выяснили, скорости протокола связи OPCua недостаточно для корректной работы механизма гриппера.  
Мы получили [результаты](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Documentation/Reports/%D0%9E%D1%82%D1%87%D0%B5%D1%82%20%D0%BF%D0%BE%20OPC%20UA%20-%20%D0%BC%D0%B5%D0%B4%D0%BB%D0%B5%D0%BD%D0%BD%D1%8B%D0%B9.pdf).  
Следующим шагом было выбор наиболее подхоящего протокола связи для достижения минимальных задержек и совместимости с нашим оборудованием:  [Таблица](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Documentation/Reports/%D0%A1%D1%80%D0%B0%D0%B2%D0%BD%D0%B5%D0%BD%D0%B8%D0%B5%20%D0%BF%D1%80%D0%BE%D1%82%D0%BE%D0%BA%D0%BE%D0%BB%D0%BE%D0%B2.pdf).  
Помимо этого мы изучали [метод внедрения СУ в Kubernetes](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Documentation/Reports/%D0%9E%D1%82%D1%87%D0%B5%D1%82%20%D0%BF%D0%BE%20kubernetes%20%D0%B8%20%D0%BF%D1%80%D0%BE%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D1%8B%20%D0%BD%D0%B0%20python%20(%D0%BF%D0%BE%D0%B4%D1%85%D0%BE%D0%B4%D0%B8%D1%82%20%D0%B4%D0%BB%D1%8F%20%D0%BB%D1%8E%D0%B1%D0%BE%D0%B3%D0%BE%20%D0%AF%D0%9F).pdf).  
####
После глубокого анализа протоколов связи, мы выяснили что EtherNetIP совместим с нашим оборудованием и удовлетворяет критерию скорости передачи данных.  
Ввиду сложности задачи, было решено протестировать EtherNetIP на двух контроллерах Siemens S7 1200, 1-ый PLC настроен как адаптер, 2-ой PLC в качестве сканнера, в качестве документации мы использовали следующие документы:
[Adapter](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Documentation/Ethernet.IP/%D0%91%D0%BE%D0%BB%D1%8C%D1%88%D0%B0%D1%8F%20%D0%BF%D0%BE%D0%B4%D1%80%D0%BE%D0%B1%D0%BD%D0%B0%D1%8F%20%D0%B4%D0%BE%D0%BA%D1%83%D0%BC%D0%B5%D0%BD%D1%82%D0%B0%D1%86%D0%B8%D1%8F/109782315_EtherNetIP_Adapter_DOC_V10_en-1.pdf) и [Scanner](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Documentation/Ethernet.IP/%D0%91%D0%BE%D0%BB%D1%8C%D1%88%D0%B0%D1%8F%20%D0%BF%D0%BE%D0%B4%D1%80%D0%BE%D0%B1%D0%BD%D0%B0%D1%8F%20%D0%B4%D0%BE%D0%BA%D1%83%D0%BC%D0%B5%D0%BD%D1%82%D0%B0%D1%86%D0%B8%D1%8F/109782314_EtherNetIP_Scanner_DOC_V1_3_en%20(1).pdf)  
####
[Проект TIA PORTAL V16](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Projects/TiaPortal_EIP_Project/eip.zap16) тестирует задержки передачи данных между двумя PLC по EtherNetIP, в ходе тестирования задержки 5-20мс.  

## Код СУ
Хотя протокол OPCua имеет некоторые проблемы при зажиме приспособления, в реальном процессе, поскольку протокол EtherNet/IP всегда имел некоторые проблемы, которые не могут быть решены, мы решили использовать протокол OPCua и язык Golang для завершения написания СУ.

Установка представляет собой набор нескольких объектов, которые выполняют свою задачу в процессе работы линии: вращающейся платформы, гриппера, упаковочной линии, сортировочной линии. 

### 1. Станция перемещения
 ![photo_2024-05-27_18-44-42](https://github.com/Spynch/HS_line316/assets/110130006/954335ad-1ccb-4438-a3f8-2222028dbce3)
 
Функции:
  - Перемещение захвата на крайнюю левую позицию над магазином фишек (там, где будет выдаваться заготовка). Захват опускается, захватывает фишку и переносит её на станцию сверления.
  - Захват забирает обработанное изделие и переносит его на станцию упаковки в центральную позицию.
  - Захват переносит коробку в крайнее правое положение и опускает коробку на конвейер.

[Код станции перемещения](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Projects/CS/hsLineOpc/internal/galaction/logic_hs.go)

### 2. Станция сверления
![photo_2024-05-27_18-44-32](https://github.com/Spynch/HS_line316/assets/110130006/494b93d8-83f6-4c23-a437-ab748e7b1125)

Функции:
  - Стол совершает поворотное движение в 4 такта и приходит к модулю проверки положения заготовки.
  - Далее идет оборот стола ещё на один такт, чтобы попасть к модулю сверления. 
  - Если заготовка является перевернутой, то модуль сверления опускается.
  - После обработки стол поворачивается к изначальной позиции, обработка детали закончена.

[Код станции сверления](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Projects/CS/hsLineOpc/internal/galaction/logic_procs.go)

### 3. Станция упаковки
![photo_2024-05-27_18-44-46](https://github.com/Spynch/HS_line316/assets/110130006/7fd1a6e3-e4ea-4f62-ad4f-d0270033e25f)

Функции:
  - Коробки поступают на позицию подготовки и открываются, после чего гриппер захватывает обработанный продукт и переносит его на центральную упаковочную станцию. Затем кладет его в коробку и упаковывает.

[Код станции упаковки](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Projects/CS/hsLineOpc/internal/galaction/logic_packs.go)

### 4. Станция распределения заготовок
 ![photo_2024-05-27_18-44-44](https://github.com/Spynch/HS_line316/assets/110130006/1e70fb35-7167-40b7-9e6b-237a5c021b98)

Функции:
  - В зависимости от определенного ранее на станции обработки цвета изделия, запускается соответствующий отсекатель.
  - коробка c заготовкой идёт в нужный накопитель.

[Код станции распределения заготовок](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Projects/CS/hsLineOpc/internal/galaction/logic_ss.go)

### 5. Имитируемый технологический процесс
 - Шаг 1: Размещение объекта на карусели
 - Шаг 2: Запуск процесса обработки
 - Шаг 3: Перемещение к упаковке
 - Шаг 4: Запуск упаковки
 - Шаг 5: Перемещение к сортировке
 - Шаг 6: Запуск сортировки
 - Шаг 7: Возврат в исходное положение

[Код Основной процесс(entry)](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Projects/CS/hsLineOpc/internal/galaction/entry.go)

## hsLineOpc
* Это [программа](https://github.com/VladislavLowPriority/CS_Packaging_line_316/tree/main/Projects/CS/hsLineOpc), которая читает управляющие теги с сервера на TS и отдает команды установке.
* Сборка:
  * Для начала надо запустить сервер на TS (см. OPCKuber+GUI)
  * Версия go 1.24 или выше:
    ```
    go version
    ```
  * Устанавливаем зависимости:
    ```
    go mod tidy
    ```
  * Ставим в .env файл свои данные (TS_SERVER_CONN будет точно другим, установки скорее всего останется таким же)
  * Запускаем программу:
    ```
    go run cmd/main.go
    ```
* Краткое описание пакетов:
  * api:
  Предоставляет функции и структуры для взаимодействия с opc клиентами
  * internal/handler:
  Оборачивает некоторые более сложные вызовы для простоты в main
  * internal/galaction:
  Код, который управляет установкой (логика ее работы)
  * main:
  Создаем клиенты opc для сервера TS и установки. Подписываемя на TS сервер (отслеживаем изменения тегов на TS сервере, а также отправляем текущие значения). Затем бесконечный цикл проверки значений с TS сервера

## OPCKuber+GUI 
### Реализация сервера OPCua и Веб-интерфейса для их размещения в контейнере Kubernetes 
В рамках стандарта Industry 4.0 и будущей модернизации программы PLC, было принято решение о реализации OPCua сервера непосредсвенно в контейнере Kubernetes. Для этого были созданы программы:

 - [Программа №1.](https://github.com/VladislavLowPriority/CS_Packaging_line_316/tree/main/Projects/OPCKuber%2BGUI/OPCkuber) Реализует OPCua сервер в Kubernetes и на данный момент синхронизирует теги в сервере с тегами на плк. Также были реализованы дополнительные теги (старт/стоп/возвращение в исходное состояние) для интерфейсного доступа к СУ с SCADA или веб-приложения.

 - [Программа №2.](https://github.com/VladislavLowPriority/CS_Packaging_line_316/tree/main/Projects/OPCKuber%2BGUI/UI_hs) Веб-интерфейс для взаимодействия с тегами старт/стоп/возвращение в исходное состояние (Клиент OPCua отправляет новое значение тегов в OPCua сервер в Kubernetes). Программа позволяет взаимодействовать с СУ, развернутой в Kubernetes.

Для запуска программ необходимо выполнить следующие шаги:  
1. Установите Node.js с официального сайта: https://nodejs.org/  
2. После установки перейдите в папку проекта и выполните:  
```bash
npm install  # Установит все зависимости из package.json
```  
3. Запуск программ:

Скрипты запуска можно можно посмотреть в файлах package.json.  
Для запуска в режиме разработки независимо от № программы можно использовать команду:  
   ```bash
   npm run dev
   ```


## Возникшие сложности:
#### Тестирование библиотек для Ethernet/IP
В ходе тестирования библиотек была выявлена проблема, что при использовании любой из них, мы не можем отправить корректный запрос в плк с указанием прочтения нужной ячейки памяти. При этом ПЛК, настроенный как EthernetIP Adapter, видится в сети как устройство EthernetIP. В мануале к библиотеке Siemens ENIP Adapter указывается конкретные Class = 0x04, Instance = 101 (Input), 102 (Output), 104 (Configuration), где должна храниться и откуда должна приниматься информация в соответствующий объект, созданный в ПЛК. Однако, при формировании любого запроса, кроме Class = 0x06, Instance = 1, который возвращает конфигурацию ПЛК (серийный номер, и тд..), ответа с ПЛК не поступает. Пример ответа: Object does not exist. При этом на сам ПЛК в поле запроса приходят корректные значения соответсвтующих атрибутов.  
Есть несколько вариантов:  
1)библиотеки СУ неправильно интерпретируют запросы  
2)библиотека в ПЛК не может общаться с устройством, отличным изготовителем от Siemens  
3)библиотека в ПЛК находится под паролем и для ее корректной работы необходим пароль

#### Аппаратные проблемы, возникшие во время отладки программы
В процессе отладки мы обнаружили, что конвейерная лента, упаковочный стол и некоторые другие части компонентов не могут работать должным образом, поэтому мы еще не завершили отладку программ, связанных с конвейерной лентой и упаковочным столом.


#### Список устройств которые не работают
<table>
  <tr>
    <th>устройства которые не работают</th>
    <th>тэги</th>
  </tr>
  <tr>
    <td>Станция сверления</td>
    <td>processing_output_2_drill_down</td>
  </tr>
  <tr>
    <td>Станция сверления</td>
    <td>processing_output_3_drill_up</td>
  </tr>
  <tr>
    <td>Станция перемещения</td>
    <td>handling_output_5_gripper_to_down</td>
  </tr>
  <tr>
    <td>Станция перемещения</td>
    <td>handling_output_6_gripper_to_open</td>
  </tr>
  <tr>
    <td>Станция перемещения</td>
    <td>handling_output_7_push_workpiece</td>
  </tr>
  <tr>
    <td>Станция упаковки</td>
    <td>packing_output_4_push_box</td>
  </tr>
  <tr>
    <td>Станция упаковки</td>
    <td>packing_output_5_fix_box_upper_side</td>
  </tr>
  <tr>
    <td>Станция упаковки</td>
    <td>packing_output_6_fix_box_tongue</td>
  </tr>
  <tr>
    <td>Станция распределения заготовок</td>
    <td>sorting_output_2_push_silver_workpiece</td>
  </tr>
  <tr>
    <td>Станция распределения заготовок</td>
    <td>sorting_output_3_push_red_workpiece</td>
  </tr>
  <tr>
    <td>Станция распределения заготовок</td>
    <td>sorting_output_1_move_conveyor_left</td>
  </tr>
  <tr>
    <td>Станция распределения заготовок</td>
    <td>sorting_output_0_move_conveyor_right</td>
  </tr>
</table>  

## Что можно сделать дальше:  
#### 1. Используя наработки и документацию EtherNet/IP настроить коммуникацию между ПЛК и ПК.  
#### 2. Взаимодействуя с командой отладки линии 316, протестировать полный цикл работы установки и программы СУ.  
#### 3. Разместить в Kubernetes программу СУ, сервер OPCUA и веб-интерфейс.  
#### 4. Заменить ПЛК Siemens S7-1200 на I/O модули, работающие по протоколу EtherNet/IP "из коробки".  
#### 5. Как альтернативу EtherNet/IP можно реализовать связь по протоколу OPC UA over TSN. Он также является протоклом реального времени, для его реализации необходим промышленный TSN коммутатор.  

## Таблица входных и выходных тэгов

### Processing station PLC
<table>
  <tr>
    <td>Тэги</td>
    <td>Описание</td>
    <td>Адрес</td>
  </tr>
  <tr>
    <th colspan="3">Входы</th>
  </tr>
  <tr>
    <td>processing_input_4_workpiece_detected</td>
    <td>заготовка обнаружена</td>
    <td>ns:4, i:5</td>
  </tr>
  <tr>
    <td>processing_input_2_workpiece_silver</td>
    <td>заготовка серебряного цвета</td>
    <td>ns:4, i:7</td>
  </tr>
  <tr>
    <td>processing_input_5_carousel_init</td>
    <td>инициализация карусели</td>
    <td>ns:4, i:3</td>
  </tr>
  <tr>
    <td>processing_input_6_hole_detected</td>
    <td>отверстие обнаружено</td>
    <td>ns:4, i:4</td>
  </tr>
  <tr>
    <td>processing_input_7_workpiece_not_black</td>
    <td>заготовка не черного цвета</td>
    <td>ns:4, i:6</td>
  </tr>
  <tr>
    <th colspan="3">Выходы</th>
  </tr>
  <tr>
    <td>processing_output_0_drill</td>
    <td>дрель</td>
    <td>ns:4, i:12</td>
  </tr>
  <tr>
    <td>processing_output_1_rotate_carousel</td>
    <td>вращение карусели</td>
    <td>ns:4, i:13</td>
  </tr>
  <tr>
    <td>processing_output_2_drill_down</td>
    <td>дрель вниз</td>
    <td>ns:4, i:14</td>
  </tr>
  <tr>
    <td>processing_output_3_drill_up</td>
    <td>дрель вверх</td>
    <td>ns:4, i:15</td>
  </tr>
  <tr>
    <td>processing_output_4_fix_workpiece</td>
    <td>фиксация заготовки</td>
    <td>ns:4, i:16</td>
  </tr>
  <tr>
    <td>processing_output_5_detect_hole</td>
    <td>детектирование отверстия</td>
    <td>ns:4, i:17</td>
  </tr>
</table>

### Handling and Packing PLC
<table>
  <tr>
    <td>Тэги</td>
    <td>Описание</td>
    <td>Адрес</td>
  </tr>
  <tr>
    <th colspan="3">Входы</th>
  </tr>
  <tr>
    <td>handling_input_0_workpiece_pushed</td>
    <td>заготовка протолкнута</td>
    <td>ns:4, i:29</td>
  </tr>
  <tr>
    <td>handling_input_1_grippe_at_right</td>
    <td>захват справа</td>
    <td>ns:4, i:32</td>
  </tr>
  <tr>
    <td>handling_input_2_gripper_at_start</td>
    <td>захват в начальном положении</td>
    <td>ns:4, i:31</td>
  </tr>
  <tr>
    <td>handling_input_3_gripper_down_pack_lvl</td>
    <td>захват опущен на уровень упаковки</td>
    <td>ns:4, i:33</td>
  </tr>
  <tr>
    <td>packing_input_7_pack_turned_on</td>
    <td>упаковка включена</td>
    <td>ns:4, i:42</td>
  </tr>
  <tr>
    <th colspan="3">Выходы</th>
  </tr>
  <tr>
    <td>handling_output_0_to_green</td>
    <td>на зеленую позицию</td>
    <td>ns:4, i:34</td>
  </tr>
  <tr>
    <td>handling_output_1_to_yellow</td>
    <td>на желтую позицию</td>
    <td>ns:4, i:35</td>
  </tr>
  <tr>
    <td>handling_output_2_to_red</td>
    <td>на красную позицию</td>
    <td>ns:4, i:36</td>
  </tr>
  <tr>
    <td>handling_output_3_gripper_to_right</td>
    <td>захват вправо</td>
    <td>ns:4, i:37</td>
  </tr>
  <tr>
    <td>handling_output_4_gripper_to_left</td>
    <td>захват влево</td>
    <td>ns:4, i:38</td>
  </tr>
  <tr>
    <td>handling_output_5_gripper_to_down</td>
    <td>захват вниз</td>
    <td>ns:4, i:39</td>
  </tr>
  <tr>
    <td>handling_output_6_gripper_to_open</td>
    <td>захват открыть</td>
    <td>ns:4, i:40</td>
  </tr>
  <tr>
    <td>handling_output_7_gripper_push_workpiece</td>
    <td>захват протолкнуть заготовку</td>
    <td>ns:4, i:41</td>
  </tr>
  <tr>
    <td>packing_output_4_push_box</td>
    <td>протолкнуть коробку</td>
    <td>ns:4, i:43</td>
  </tr>
  <tr>
    <td>packing_output_5_fix_box_upper_side</td>
    <td>фиксация верхней стороны коробки</td>
    <td>ns:4, i:44</td>
  </tr>
  <tr>
    <td>packing_output_6_fix_box_tongue</td>
    <td>фиксация язычка коробки</td>
    <td>ns:4, i:45</td>
  </tr>
  <tr>
    <td>packing_output_7_pack_box</td>
    <td>упаковать коробку</td>
    <td>ns:4, i:46</td>
  </tr>
</table>

### Sorting station PLC
<table>
  <tr>
    <td>Тэги</td>
    <td>Описание</td>
    <td>Адрес</td>
  </tr>
  <tr>
    <th colspan="3">Входы</th>
  </tr>
  <tr>
    <td>sorting_input_3_box_on_conveyor</td>
    <td>коробка на конвейере</td>
    <td>ns:4, i:9</td>
  </tr>
  <tr>
    <td>sorting_input_4_box_is_down</td>
    <td>коробка опущена</td>
    <td>ns:4, i:10</td>
  </tr>
  <tr>
    <th colspan="3">Выходы</th>
  </tr>
  <tr>
    <td>sorting_output_0_move_conveyor_right</td>
    <td>движение конвейера вправо</td>
    <td>ns:4, i:19</td>
  </tr>
  <tr>
    <td>sorting_output_1_move_conveyor_left</td>
    <td>движение конвейера влево</td>
    <td>ns:4, i:20</td>
  </tr>
  <tr>
    <td>sorting_output_2_push_silver_workpiece</td>
    <td>протолкнуть серебряную заготовку</td>
    <td>ns:4, i:21</td>
  </tr>
  <tr>
    <td>sorting_output_3_push_red_workpiece</td>
    <td>протолкнуть красную заготовку</td>
    <td>ns:4, i:22</td>
  </tr>
</table>

### Управляющие сигналы (OPC UA сервер Kubernetes)
#### Точно такие же теги заведены в OPC UA сервер Kubernetes, меняется только namespace = 1
<table>
  <tr>
    <td>Тэги</td>
    <td>Описание</td>
    <td>Адрес</td>
  </tr>
  <tr>
    <td>Start_hs</td>
    <td>старт</td>
    <td>ns:1, i:51</td>
  </tr>
  <tr>
    <td>stop_hs</td>
    <td>стоп</td>
    <td>ns:1, i:52</td>
  </tr>
  <tr>
    <td>back_to_start</td>
    <td>возвращение в исходное состояние</td>
    <td>ns:1, i:53</td>
  </tr>
</table>
