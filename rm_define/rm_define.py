# Robot mode "enum".
robot_mode_gimbal_follow = 1
robot_mode_chassis_follow = 2
robot_mode_free = 3

# Objective: Sets the travel mode
#
# * Chassis Lead Mode: The gimbal follows the chassis to rotate along the yaw
#   axis.
# * Gimbal Lead Mode: The chassis follows the gimbal to rotate along the yaw
#   axis.
# * Free Mode: The gimbal and the chassis move without affecting each other.
def robot_set_mode(mode_enum):
    # Intentionally do nothing for now. We could potentially keep track of
    # robot state in case we want to do some actual real similation of it.
    print(f'rm_define.robot_set_mode({mode_enum})')
