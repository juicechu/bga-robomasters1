# If true, sweep from side to side while looking for targets.
TARGET_SEARCHING  = True

# Angle of sweep (90 means -45 to 45).
TARGET_SEARCHING_ANGLE = 90

# Type of target to track.
# 0: Robot
# 1: Vision Marker
TARGET_TYPE = 0

# Method to use for gimbal movement when tracking.
# 0: Gimbal speed.
# 1: Angle to target.
TARGET_TRACKING_MODE = 0

# If true, automatically fire on lock.
AUTO_FIRE_ON_LOCK = True

# Maximum distance (in meters) the target must be to be fired upon when
# AUTO_FIRE_ON_LOCK is enabled.
AUTO_FIRE_MAX_DISTANCE = 2.0

# If true, uses PID for gimbal movement.
PID_ENABLED = True

# PID controller setup ([P, I, D]).
PID_PITCH_PARAMETERS = [0, 0, 0]
PID_YAW_PARAMETERS   = [0, 0, 0]

# If true, allows controller override (i.e. manual control) of the Robot when
# Sentry Mode is enabled.
CONTROLLER_OVERRIDE = False

# Program entry point. Set up robot and start looking for targets.
def start():
    # Enable S1 robot identification.
    vision_ctrl.enable_detection(rm_define.vision_detection_car)

    # Set reasonable gimbal speed for finding new targets.
    gimbal_ctrl.set_rotate_speed(60)

    # Move gimbal and chassis independently.
    robot_ctrl.set_mode(rm_define.robot_mode_free)

    # Rotating white leds (searching for targets).
    led_ctrl.set_top_led(rm_define.armor_top_all, 255, 255, 255,
            rm_define.effect_marquee)

    if TARGET_SEARCHING:
        while True:
	    # Sweep from side to side.
            media_ctrl.play_sound(rm_define.media_sound_gimbal_rotate)
            gimbal_ctrl.yaw_ctrl(-90)
            media_ctrl.play_sound(rm_define.media_sound_gimbal_rotate)
            gimbal_ctrl.yaw_ctrl(90)
    else:
        while True:
            # Just sleep as there is nothing to do other than waiting for
            # a target to be identified.
            time.sleep(60)

# Return the bounding box information (X, Y, W, H) for the closest detected
# robot.
def FindClosestRobot(robotDetectionInfo):
    return FindClosest(robotDetectionInfo, 4)

# Simple algorithm to find the closest detected object. It simply iterates
# through all of the detected objects and returns the one with the bigger
# height.
#
# Other than the detection info, it takes as parameter the expected number of
# items per detected object as it varies for a detected robot or a detected
# vision marker, for example.
def FindClosest(detectionInfo, numEntriesPerObject):
    numEntries = len(detectionInfo) - 1  # Ignore size entry.
    numObjects = numEntries // numEntriesPerObject
    if numObjects != detectionInfo[0]:
        # Got an unexpected number of entries.
        return None

    modulo = numEntries % numEntriesPerObject
    if modulo != 0:
        # Got incomplete number of entries.
        return None

    closestHeight = 0.0  # Impossible height.
    closestIndex = 1 # Defaults to first Robot detected.

    # Check height of the bounding box of each detected robot. Return the one
    # with the biggest height.
    for i in range(1, len(detectionInfo) - 1, numEntriesPerObject):
        objectHeight = detectionInfo[i + numEntriesPerObject - 1]
        if objectHeight > closestHeight:
            # Found a bigger height.
            closestHeight = objectHeight
            closestIndex = i

    # Return only the relevant info about the selected object.
    return detectionInfo[closestIndex:closestIndex + numEntriesPerObject]

# Known height of a Robomaster S1 in millimeters and inches.
ROBOT_KNOWN_HEIGHT_MM = 270.0
ROBOT_KNOWN_HEIGHT_IN = 10.6

# Due to the coordinate system used, the focal length can be inferred directly
# (you can still compute it yourself using FocalLength() to see it matches).
ROBOT_CAMERA_FOCAL_LENGTH = 1.0

# Returns the distance to the object in the same unit as the given
# knownHeightOrWidth.
def Distance(knownHeightOrWidth, focalLength, heightOrWidth):
    return (knownHeightOrWidth * focalLength) / heightOrWidth

# Compute distance in millimeters to a detected Robomaster S1 given its
# bounding box height.
def DistanceToRobotMM(height):
    return Distance(ROBOT_KNOWN_HEIGHT_MM, ROBOT_CAMERA_FOCAL_LENGTH,
            height)

# Robomaster S1 camera field of view information.
ROBOT_CAMERA_HORIZONTAL_FOV = 96
ROBOT_CAMERA_VERTICAL_FOV   = 54

# Aim return codes.
AIM_ERROR       = 0
AIM_IN_PROGRESS = 1
AIM_DONE        = 2

# Aims the Robomaster S1 gimbal to the given coordinates.
#
# There are 2 possible operating modes supported:
#
# Direct Mode computes the angle between the source (sight) and destination
# (detected object) positions and uses that angle to directly move the gimbal to
# the target. In this mode only dst_x and dst_y should be provided.
#
# PID Mode computes the delta between the source (sight) and destination
# (detected object) positions and feeds this delta as errors to the given PID
# controllers. The gimbal turn speed is then set based on the PID controllers
# output.
#
# Return values:
#
# AIM_ERROR indicates invalid parameters were provided.
# AIM_IN_PROGRESS indicates we are currently trying to aim at the destination
#   position but did not get a lock yet.
# AIM_DONE indicates we are now locked to the destination position.
#
# Note that no matter which method is used, at least 2 passes are required to
# get an AIM_DONE as the first pass will always return AIM_IN_PROGRESS (unless
# the robot was already pointing directly to the destination position).
#
# Actions that require a target lock should only be done when AIM_DONE is
# returned.
def Aim(dst_x, dst_y, pid_yaw = None, pid_pitch = None):
    if dst_x < 0.0 or dst_x > 1.0 or dst_y < 0.0 or dst_y > 1.0:
        # Invalid dst_x or dst_y.
        return AIM_ERROR

    if ((pid_yaw is not None and pid_pitch is None) or
            (pid_pitch is not None and pid_yaw is None)):
        # Only one of pid_yaw and pid_pitch was provided.
        return AIM_ERROR

    # Obtain sight position. This takes into account sight calibration.
    sight_info = media_ctrl.get_sight_bead_position()
    src_x = sight_info[0]
    src_y = sight_info[1]

    # Compute deltas between source and destination.
    delta_x = dst_x - src_x
    delta_y = src_y - dst_y

    if abs(delta_x) <= 0.1 and abs(delta_y) <= 0.1:
        # We are centered in the target already. There is nothing else to do.
        if pid_yaw is not None:
            # We are in PID mode. Stop gimbal rotation that might still be in
            # progress.
            gimbal_ctrl.rotate_with_speed(0, 0)
        return AIM_DONE

    if pid_yaw is not None:
        # PID mode.

        # Set error in the PID controllers.
        pid_yaw.set_error(delta_x)
        pid_pitch.set_error(delta_y)

        # Set gimbal rotation speed based on PID controllers output.
        gimbal_ctrl.rotate_with_speed(pid_yaw.get_output(),
                pid_pitch.get_output())
    else:
        # Direct mode.

        # Get current gimbal yaw and pitch angles.
        gimbal_yaw_angle = gimbal_ctrl.get_axis_angle(rm_define.gimbal_axis_yaw)
        gimbal_pitch_angle = gimbal_ctrl.get_axis_angle(
                rm_define.gimbal_axis_pitch)

        # Compute deltas between source and destination angles.
        delta_yaw_angle = ROBOT_CAMERA_HORIZONTAL_FOV * (delta_x)
        delta_pitch_angle = ROBOT_CAMERA_VERTICAL_FOV * (delta_y)

        print(gimbal_yaw_angle + delta_yaw_angle)
        
        # Move gimbal so the sight points directly to the target.
        gimbal_ctrl.angle_ctrl(gimbal_yaw_angle + delta_yaw_angle,
                gimbal_pitch_angle + delta_pitch_angle)

    return AIM_IN_PROGRESS

def vision_recognized_car(msg):
    # Create PID controllers for pitch and yaw.
    pid_pitch = rm_ctrl.PIDCtrl()
    pid_yaw = rm_ctrl.PIDCtrl()

    # Set contoller parameters.
    pid_pitch.set_ctrl_params(90,0,3)
    pid_yaw.set_ctrl_params(120,0,5)
    
    # Keep track of previous aim status.
    previous_aim_status = AIM_ERROR

    while True:
        robot_detection_info = vision_ctrl.get_car_detection_info()
        if robot_detection_info[0] == 0:
            break

        print(f'Seeing {robot_detection_info[0]} robots.')

        closest_robot_info = FindClosestRobot(robot_detection_info)
        if closest_robot_info is None:
            print(f'Unexpected robot data. Abort tracking.')
            break
            
        distance = DistanceToRobotMM(closest_robot_info[3])
        if distance is None:
            print(f'Can\'t get distance. Abort tracking.')
            break

        distance_in_meters = distance / 1000

        print(f'Closest robot is {distance_in_meters:.2f} meters away.')

        aim_status = Aim(closest_robot_info[0], closest_robot_info[1],
                         pid_yaw, pid_pitch)
        if aim_status == AIM_DONE:
            print('Target locked.')

            if previous_aim_status != aim_status:
                # Rotating red lights (Target locked).
                led_ctrl.set_top_led(rm_define.armor_top_all, 255, 0, 0,
                    rm_define.effect_marquee)

            if distance_in_meters <= 2.0:
		print(f'Fire!')
            	gun_ctrl.fire_once()
            else:
                print(f'Too far. Not firing.')

        else:
            if aim_status == AIM_IN_PROGRESS:
                print(f'Aiming...')

                if previous_aim_status != aim_status:
                    # Rotating yellow lights (tracking target).
                    led_ctrl.set_top_led(rm_define.armor_top_all, 255, 255, 0,
                        rm_define.effect_marquee)

		# Give some time for the gimbal position to stabilize as
                # otherwise we might get bogus target position data.
                time.sleep(0.1)

        previous_aim_status = aim_status
    
    gimbal_ctrl.rotate_with_speed(0, 0)

    # Back to rotating white leds (searching for targets).
    led_ctrl.set_top_led(rm_define.armor_top_all, 255, 255, 255,
            rm_define.effect_marquee)
