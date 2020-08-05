export function table_template = ({id, startdate, track, weather_randomness, P_hourOfDay, P_timeMultiplier, P_sessionDurationMinute,  Q_hourOfDay, Q_timeMultiplier, Q_sessionDurationMinute, R_hourOfDay, R_timeMultiplier, R_sessionDurationMinute, pit_window_length_sec, is_refuelling_allowed_in_race, mandatory_pitstop_count, is_mandatory_pitstop_refuelling_required, is_mandatory_pitstop_tyre_change_required, is_mandatory_pitstop_swap_driver_required, tyre_set_count}) => {
    return 
    `<table><tr>
    <td>ID</td><td>${id}</td>
    </tr>
    <tr>
    <td>Start Date</td><td>${startdate}</td>
    </tr>
    <tr>
    <td>Track</td><td>${track}</td>
    </tr>
    <tr>
    <td>Weather Randomness</td><td>${weather_randomness}</td>
    </tr>
    <tr>
    <td>Practice Start Hour of Day</td><td>${P_hourOfDay}</td>
    </tr>
    <tr>
    <td>Practice Time Multiplier</td><td>&times;${P_timeMultiplier}</td>
    </tr>
    <tr>
    <td>Practice Session Duration</td><td>${P_sessionDurationMinute} min.</td>
    </tr>
    <td>Qualifying Start Hour of Day</td><td>${Q_hourOfDay}</td>
    </tr>
    <tr>
    <td>Qualifying Time Multiplier</td><td>&times;${Q_timeMultiplier}</td>
    </tr>
    <tr>
    <td>Qualifying Session Duration</td><td>${Q_sessionDurationMinute} min.</td>
    </tr>
    <td>Race Start Hour of Day</td><td>${R_hourOfDay}</td>
    </tr>
    <tr>
    <td>Race Time Multiplier</td><td>&times;${R_timeMultiplier}</td>
    </tr>
    <tr>
    <td>Race Session Duration</td><td>${R_sessionDurationMinute} min.</td>
    </tr>
    <tr>
    <td>Pit Window Length</td><td>${pit_window_length_sec} sec. (${pit_window_length_sec / 60} min.)</td>
    </tr>
    <tr>
    <td>Refuelling Allowed</td><td>${is_refuelling_allowed_in_race == false ? "NO" : "YES"}</td>
    </tr>
    <tr>
    <td>Mandatory Pitstop Count</td><td>${mandatory_pitstop_count}</td>
    </tr>
    <tr>
    <td>Refuelling Required in Mandatory Pitstop?</td><td>${is_mandatory_pitstop_refuelling_required== false ? "NO" : "YES"}</td>
    </tr>
    <tr>
    <td>Tyre Change Required in Mandatory Pitstop?</td><td>${is_mandatory_pitstop_tyre_change_required == false ? "NO" : "YES"}</td>
    </tr>
    <tr>
    <td>Driver Swap Required in Mandatory Pitstop?</td><td>${is_mandatory_pitstop_swap_driver_required == false ? "NO" : "YES"}</td>
    </tr>
    <tr>
    <td>Tyre Set Count</td><td>${tyre_set_count}</td>
    </tr></table>`;
