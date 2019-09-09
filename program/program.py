import rm_define.rm_define as rm_define

# Sample usage of the Robomaster S1 stubs. start() is called automatically.
def start():
    # Set robot mode. Note this call is identical to the one that would be
    # done in the actual Robomaster S1 App.
    rm_define.robot_set_mode(rm_define.robot_mode_free)
