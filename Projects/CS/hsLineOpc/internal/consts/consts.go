package consts

// TS
const (
	TsStartHs     = "ns=1;i=51"
	TsStopHs      = "ns=1;i=52"
	TsBackToStart = "ns=1;i=53"
)

// inputs
const (
	InputBoxOnConveyor = "ns=4;i=9"
	InputBoxIsDown     = "ns=4;i=10"
	InputRed           = "ns=4;i=24"
	InputBlc           = "ns=4;i=25"
	InputSil           = "ns=4;i=26"

	InputCarouselRotation = "ns=4;i=3"
	InputM5               = "ns=4;i=4"
	InputRedAndSil        = "ns=4;i=6"
	InputSil2             = "ns=4;i=7"

	InputGripperStart    = "ns=4;i=31"
	InputGripperPack     = "ns=4;i=30"
	InputGripperConveyor = "ns=4;i=32"
)

// outputs
const (
	OutputConveyorRight = "ns=4;i=19"
	OutputConveyorLeft  = "ns=4;i=20"

	OutputGripperLeft   = "ns=4;i=38"
	OutputGripperRight  = "ns=4;i=37"
	OutputGripperOpen   = "ns=4;i=40"
	OutputGripperUpDown = "ns=4;i=39"
)

// hs
// # input tags
// gripper_start_sensor = 'ns=4;i=31'
// gripper_pack_sensor = 'ns=4;i=30'
// gripper_conveyor_sensor = 'ns=4;i=32'
// -------------
// # output tags
// gripper_toggle_up_down = 'ns=4;i=39'
// gripper_open = 'ns=4;i=40'
// gripper_move_left = 'ns=4;i=38'
// gripper_move_right = 'ns=4;i=37'
// drop_puck = 'ns=4;i=41'
// green_tag = 'ns=4;i=34'
// yellow_tag = 'ns=4;i=35'

// packs
// # output tags
// fix_box_upper_side = 'ns=4;i=44'
// fix_box_tongue = 'ns=4;i=45'
// pack_box = 'ns=4;i=46'

// procs
// # Color tags:
// red_tag = 'ns=4;i=24'
// silver_tag = 'ns=4;i=26'
// black_tag = 'ns=4;i=25'
// -------------
// # Input tags:
// carousel_rotation_tag = 'ns=4;i=3'
// m5_tag = 'ns=4;i=4'
// red_and_silvery = 'ns=4;i=6'
// silvery = 'ns=4;i=7'
// -------------
// # Output tags:
// drilling = "ns=4;i=12"
// carousel_rotate = "ns=4;i=13"
// drill_down = "ns=4;i=14"
// drill_up = "ns=4;i=15"
// m4_toggle = 'ns=4;i=16'
// m5_toggle = 'ns=4;i=17'

// ss
// # input
// box_on_conveyor_tag = 'ns=4;i=9'
// box_is_down_tag = 'ns=4;i=10'
// red_tag = 'ns=4;i=24'
// black_tag = 'ns=4;i=25'
// silver_tag = 'ns=4;i=26'
// ----------------
// # output
// move_conveyor_right = 'ns=4;i=19'
// move_conveyor_left = 'ns=4;i=20'
// push_silver = 'ns=4;i=21'
// push_red = 'ns=4;i=22'
