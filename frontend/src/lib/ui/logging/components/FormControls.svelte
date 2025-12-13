<script lang="ts">
    import {qsoState} from "$lib/states/new-qso-state.svelte";
    import {handleAsyncError} from "$lib/utils/error-handler";
    import {LogQso, CurrentSessionQsoSlice, TotalQsosByLogbookId} from "$lib/wailsjs/go/facade/Service";
    import {types} from "$lib/wailsjs/go/models";
    import {configState} from "$lib/states/config-state.svelte";
    import {showToast} from "$lib/utils/toast";
    import {isValidCallsignForLog} from "$lib/constants/callsign";
    import {sessionState} from "$lib/states/session-state.svelte";
    import {isContestMode} from "$lib/stores/logging-mode-store";
    import {contestTimers} from "$lib/utils/contest-timers.svelte";
    import {contestState} from "$lib/states/contest-state.svelte";
    import {catState} from "$lib/states/cat-state.svelte";
    import {shortcut} from "@svelte-put/shortcut";

    let isLogging: boolean = $state(false);

    // We must calculate the power value here, because it depends on the panel which displays the value
    // is not guaranteed to be loaded into the DOM.
    const calculateTxPwr = (): number => {
        let pwr = parseInt(catState.txPower);
        if (isNaN(pwr)) {
            pwr = configState.default_power;
        }
        if (configState.use_power_multiplier) {
            pwr = pwr * configState.power_multiplier;
        }
        return pwr;
    }

    const resetAction = (): void => {
        qsoState.stopTimer();
        qsoState.reset();
        isLogging = false;
        if ($isContestMode) {
            const srxElem = document.getElementById('srx_rcvd') as HTMLInputElement;
            if (srxElem) {
                srxElem.classList.remove('outline-red-500', 'outline-2');
                srxElem.classList.add('outline-gray-300', 'outline-1');
            }
        }
        const elem = document.getElementById('call') as HTMLInputElement;
        if (elem) elem.focus();
    }

    const canLog = (): boolean => {
        return isValidCallsignForLog(qsoState.call)
    };

    const logContact = async (): Promise<void> => {
        if (!canLog()) {
            const elem = document.getElementById('call') as HTMLInputElement;
            if (elem) {
                elem.focus();
                elem.select();
            }
            return;
        }
        if (isLogging) return; // Prevent double-clicks
        isLogging = true;
        try {
            qsoState.tx_pwr = calculateTxPwr().toString();
            const qso: types.Qso = qsoState.toQso();
            qso.logbook_id = configState.logbook.id;

            if ($isContestMode){
                contestTimers.reset();
                contestState.totalQsos = await TotalQsosByLogbookId(configState.logbook.id);
                const stxElem = document.getElementById('stx_sent') as HTMLInputElement;
                if (stxElem) {
                    qso.stx = stxElem.value;
                    // This will auto-update the stx_sent field in the UI (see QsoPanel.svelte)
                    contestState.increment();
                }
            }

            await LogQso(qso);
            showToast.SUCCESS("QSO logged.");
            sessionState.update(await CurrentSessionQsoSlice());


        } catch (e: unknown) {
            handleAsyncError(e, 'FormControls.svelte->logContact()');
        }

        resetAction();
    }
</script>

<div class="flex w-[230px] justify-end gap-x-3">
    <button
            id="log-contact-btn"
            onclick={logContact}
            type="button"
            disabled={!canLog}
            class="disabled:bg-gray-400 disabled:cursor-not-allowed h-9 cursor-pointer rounded-md bg-indigo-600 p-2.5 py-1.5 text-base font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            title="Ctrl-s">Log Contact
    </button>
    <button
            onclick={resetAction}
            type="button"
            class="h-9 w-[74px] cursor-pointer rounded-md bg-white px-2.5 py-1.5 text-base font-semibold ring-1 shadow-sm ring-gray-300 ring-inset hover:bg-gray-100"
            title="ESC">Clear
    </button>
</div>
<svelte:window
        use:shortcut={{
        trigger: [
            {key: 'Escape', callback: resetAction},
            {key: 's', modifier: 'ctrl', callback: logContact},
        ],
    }}
></svelte:window>