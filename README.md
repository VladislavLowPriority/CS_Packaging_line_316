# Упаковочная линия 316

# Документация проекта
## Авторы
Качалов Владислав  
Полетаев Александр  
Соловьев Михаил  
Васенин Константин  
Чжан Цзяпэн  
Ван Цзыи  
Руководитель: **Куклин Егор Вадимович**  

## Описание поставленной задачи
Обеспечить сокращение задержек передачи управляющих сигналов между установкой упаковочной станции (СУ 316) и системы управления, развернутой на сервере, путем поиска и интегрирования протокола реального времени (RT). Также необходимо определить подходящий язык программирования и адаптировать под него логику СУ, подготовить интерфейс для реализации функционала старта/стопа и развернуть СУ в Kubernetes.

## Используемое ПО
Siemens TIA Portal V16 / SIMATIC STEP 7 incl. Safety and WinCC V16    
UaExpert  
VisualCode  
EnIPExploler  
Docker Desktop  
Kubernetes  


## Ссылка на проект
[HS Line 316](https://drive.google.com/drive/folders/10Y2lL00LheItdtCxrIDzM5lAGOLko-ze?usp=sharing)

## Архитектура
ПЛК взаимодействут друг с другом благодаря стандарту Profinet 
![изображение](https://github.com/n0th1ngs89/HS_Line_316_I-O/assets/146949002/05848d01-be59-402f-ab1c-b9a10d6a265b)

## Ход работы:
Была написана тестовая программа СУ на языке Golang, используюящая протокол связи OPCua.  
В ходе тестирования мы выяснили, скорости протокола связи OPCua недостаточно для корректной работы механизма гриппера.  
Мы получили такие [результаты](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Documentation/Reports/%D0%9E%D1%82%D1%87%D0%B5%D1%82%20%D0%BF%D0%BE%20OPC%20UA%20-%20%D0%BC%D0%B5%D0%B4%D0%BB%D0%B5%D0%BD%D0%BD%D1%8B%D0%B9.pdf)  
Следующим шагом было выбор наиболее подхоящего протокола связи для достижения минимальных задержек и совместимости с нашим оборудованием:  
[Таблица](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Documentation/Reports/%D0%A1%D1%80%D0%B0%D0%B2%D0%BD%D0%B5%D0%BD%D0%B8%D0%B5%20%D0%BF%D1%80%D0%BE%D1%82%D0%BE%D0%BA%D0%BE%D0%BB%D0%BE%D0%B2.pdf)  
Помимо этого мы изучали [метод внедрения СУ в Kubernetes](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Documentation/Reports/%D0%9E%D1%82%D1%87%D0%B5%D1%82%20%D0%BF%D0%BE%20kubernetes%20%D0%B8%20%D0%BF%D1%80%D0%BE%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D1%8B%20%D0%BD%D0%B0%20python%20(%D0%BF%D0%BE%D0%B4%D1%85%D0%BE%D0%B4%D0%B8%D1%82%20%D0%B4%D0%BB%D1%8F%20%D0%BB%D1%8E%D0%B1%D0%BE%D0%B3%D0%BE%20%D0%AF%D0%9F).pdf)  
##
После глубокого анализа протоколов связи, мы выяснили что EtherNetIP совместим с нашим оборудованием и удовлетворяет критерию скорости передачи данных.  
В веду сложности задачи, было решено протестировать EtherNetIP на двух контроллерах Siemens S7 1200, 1-ый PLC настроен как адаптер, 2-ой PLC в качестве сканнера, в качестве документации мы использовали следующие документы:
[адаптер](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Documentation/Ethernet.IP/%D0%91%D0%BE%D0%BB%D1%8C%D1%88%D0%B0%D1%8F%20%D0%BF%D0%BE%D0%B4%D1%80%D0%BE%D0%B1%D0%BD%D0%B0%D1%8F%20%D0%B4%D0%BE%D0%BA%D1%83%D0%BC%D0%B5%D0%BD%D1%82%D0%B0%D1%86%D0%B8%D1%8F/109782315_EtherNetIP_Adapter_DOC_V10_en-1.pdf)
и  
[сканнер](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Documentation/Ethernet.IP/%D0%91%D0%BE%D0%BB%D1%8C%D1%88%D0%B0%D1%8F%20%D0%BF%D0%BE%D0%B4%D1%80%D0%BE%D0%B1%D0%BD%D0%B0%D1%8F%20%D0%B4%D0%BE%D0%BA%D1%83%D0%BC%D0%B5%D0%BD%D1%82%D0%B0%D1%86%D0%B8%D1%8F/109782314_EtherNetIP_Scanner_DOC_V1_3_en%20(1).pdf)  
##
Проект TIA PORTAL V16 тестирует задержки передачи данных между двумя PLC по EtherNetIP, в ходе тестирования задержки 5-20мс.  
[проект](https://github.com/VladislavLowPriority/CS_Packaging_line_316/blob/main/Documentation/eip.zap16)
##
В рамках стандарта Industry 4.0 и будущей модернизации программы PLC, было принято решение о реализации OPCua сервера непосредсвенно в контейнере Kubernetes. Для этого были созданы программы:

[Программа №1.](https://github.com/VladislavLowPriority/CS_Packaging_line_316/tree/main/%D0%92%D0%BB%D0%B0%D0%B4%D0%B8%D1%81%D0%BB%D0%B0%D0%B2/OPCkuberAndGUI/OPCkuber) Реализует OPCua сервер в Kubernetes и на данный момент синхронизирует теги в сервере с тегами на плк. Также были реализованы дополнительные теги (старт/стоп/возвращение в исходное состояние) для интерфейсного доступа к СУ с SCADA или веб-приложения.

[Программа №2.](https://github.com/VladislavLowPriority/CS_Packaging_line_316/tree/main/%D0%92%D0%BB%D0%B0%D0%B4%D0%B8%D1%81%D0%BB%D0%B0%D0%B2/OPCkuberAndGUI/UI_hs) Веб-интерфейс для взаимодействия с тегами старт/стоп/возвращение в исходное состояние (Клиент OPCua отправляет новое значение тегов в OPCua сервер в Kubernetes). Программа позволяет взаимодействовать с СУ, развернутой в Kubernetes.
## Возникшие сложности:
#### Тестирование библиотек для Ethernet/IP
В ходе тестирования библиотек была выявлена проблема, что при использовании любой из них, мы не можем отправить корректный запрос в плк с указанием прочтения нужной ячейки памяти. При этом ПЛК, настроенный как EthernetIP Adapter, видится в сети как устройство EthernetIP. В мануале к библиотеке Siemens ENIP Adapter указывается конкретные Class = 0x04, Instance = 101 (Input), 102 (Output), 104 (Configuration), где должна храниться и откуда должна приниматься информация в соответствующий объект, созданный в ПЛК. Однако, при формировании любого запроса, кроме Class = 0x06, Instance = 1, который возвращает конфигурацию ПЛК (серийный номер, и тд..), ответа с ПЛК не поступает. Пример ответа: Object does not exist. При этом на сам ПЛК в поле запроса приходят корректные значения соответсвтующих атрибутов.  
Есть несколько вариантов:  
1)библиотеки СУ неправильно интерпретируют запросы  
2)библиотека в ПЛК не может общаться с устройством, отличным изготовителем от Siemens  
3)библиотека в ПЛК находится под паролем и для ее корректной работы необходим пароль

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

