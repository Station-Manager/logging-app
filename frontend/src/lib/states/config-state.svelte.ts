import { types } from '$lib/wailsjs/go/models';

export interface ConfigState {
    default_freq: string;
    default_mode: string;
    logbook: types.Logbook;
    owners_callsign: string;
    rig_name: string;
    load(this: ConfigState, cfg: types.UiConfig): void;
}
export const configState: ConfigState = $state({
    default_freq: '',
    default_mode: '',
    logbook: new types.Logbook(),
    owners_callsign: '',
    rig_name: '',

    load(this: ConfigState, cfg: types.UiConfig): void {
        this.default_freq = cfg.default_freq;
        this.default_mode = cfg.default_mode;
        this.logbook = cfg.logbook;
        this.owners_callsign = cfg.owners_callsign;
        this.rig_name = cfg.rig_name;
    },
});
