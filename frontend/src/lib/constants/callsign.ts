export const CALLSIGN_MIN = 3;
export const CALLSIGN_MAX = 12;
export const CALLSIGN_PATTERN = new RegExp(
    `^(?!\\/)(?!.*\\/$)(?!.*\\/\\/)(?=.*[A-Z])(?=.*\\d)[A-Z0-9/]{${CALLSIGN_MIN},${CALLSIGN_MAX}}$`
);
