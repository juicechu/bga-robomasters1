# Manual control of firing rate.
#
# This allows one to define a firing rate and to also enable simulated fire
# by using the trajectory light.
class FireControl:
    def __init__(self, max_rate_per_second, simulated):
        self.max_rate_per_second = max_rate_per_second
        self.simulated = simulated

        self.last_fire = time.time()

    def fire:
        delta_time = time.time() - self.last_fire
        if delta_time < 1 / (self.max_rate_per_second):
            # We can not fire now.
            return

        self.last_fire = time.time()

        if self.simulated:
            # Simulte firing with the trajectory light.
            led_ctrl.gun_on()
            led_ctrl_gun_off()
        else:
            # Fire blaster. 
            gun_ctrl.fire_once()


