<script lang="ts">
    import {qsoState} from "$lib/states/new-qso-state.svelte";
    import {handleAsyncError} from "$lib/utils/error-handler";
    import {CurrentSessionQsoSlice, LogQso, TotalQsosByLogbookId} from "$lib/wailsjs/go/facade/Service";
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
    import {clipboardState} from "$lib/states/clipboard-state.svelte";

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

    /**
     * Resets the current QSO state and associated UI elements.
     *
     * This function is responsible for stopping the QSO timer, resetting the state,
     * and updating the UI to its initial state. Additionally, it handles specific resets
     * for contest mode by modifying class styles of relevant elements. Once the reset
     * is complete, it ensures focus is set to the primary input field.
     *
     * Effects:
     * - Stops the QSO timer.
     * - Resets the QSO state to its initial values.
     * - Updates the logging state to inactive.
     * - Adjusts the styles of the contest-related input field if contest mode is active.
     * - Sets focus to the input field for entering callsigns.
     */
    const resetAction = (): void => {
        qsoState.stopTimer();
        qsoState.reset();
        isLogging = false;
        if ($isContestMode) {
            const srxElem = document.getElementById('srx_rcvd') as HTMLInputElement;
            if (srxElem) {
                srxElem.classList.remove('outline-red-500', 'outline-2', 'focus:outline-red-600');
                srxElem.classList.add('outline-gray-300', 'outline-1', 'focus:outline-indigo-600');
            }
        }
        const elem = document.getElementById('call') as HTMLInputElement;
        if (elem) elem.focus();
    }

    const canLog = (): boolean => {
        return isValidCallsignForLog(qsoState.call)
    };

    /**
     * Updates the QSO (contact log) state with the necessary values derived from
     * configuration, session, and contest states. Handles validation and state adjustments
     * when the application is in contest mode.
     *
     * Functionality:
     * - Updates the `tx_pwr` property of `qsoState` using a calculated transmission power value.
     * - Sets the `operator` property of `qsoState` based on the current session's operator value
     *   if it is non-empty.
     * - Assigns the owner's callsign to the `owner_callsign` property of `qsoState`.
     * - In contest mode:
     *   - Validates the "SRX received" field in the UI and shows an error message if invalid,
     *     along with updating input styles to indicate the error.
     *   - Sets the `stx` property of `qsoState` using the "STX sent" field value.
     *   - Automatically updates the contest's STX counter if logging is successful.
     *
     * Error Handling:
     * - Displays an error toast and updates the UI input field styling when the SRX value in
     *   contest mode is invalid.
     *
     * Side Effects:
     * - Modifies `qsoState` with updated QSO-related information.
     * - May stop the logging process (`isLogging = false`) if validation fails in contest mode.
     * - Updates UI elements and contest state in contest mode based on successful data processing.
     *
     * Parameters: None
     *
     * Returns: None
     */
    const updateQsoState = (): void => {
        qsoState.tx_pwr = calculateTxPwr().toString();
        if (sessionState.operator.trim() !== "") {
            qsoState.operator = sessionState.operator;
        }
        qsoState.owner_callsign = configState.owners_callsign;

        if ($isContestMode) {
            const srxElem = document.getElementById('srx_rcvd') as HTMLInputElement;
            if (srxElem) {
                if (srxElem.value.length < 1) {
                    showToast.ERROR("Invalid SRX value.");
                    srxElem.classList.remove('outline-gray-300', 'outline-1', 'focus:outline-indigo-600');
                    srxElem.classList.add('outline-red-500', 'outline-2', 'focus:outline-red-600');
                    srxElem.focus();
                    isLogging = false;
                    return;
                }
            }
            // Increment only when we know we are able to log the QSO
            const stxElem = document.getElementById('stx_sent') as HTMLInputElement;
            if (stxElem) {
                qsoState.stx = stxElem.value;
                // This will auto-update the stx_sent field in the UI (see QsoPanel.svelte)
                contestState.increment();
            }
        }
    }

    /**
     * Asynchronously logs a QSO (contact) entry into the logbook. Handles various
     * states and errors during the process to ensure a consistent logging operation.
     *
     * Function Behavior:
     * - If logging is not allowed (`canLog()` returns false), focuses and selects the
     *   "call" input field and exits.
     * - Prevents duplicate log attempts by checking and setting a locking mechanism (`isLogging`).
     * - Updates the QSO state and retrieves the QSO object to be logged.
     * - Assigns the current logbook id to the QSO and performs the log operation via `LogQso`.
     * - Displays a success message upon successful logging and updates necessary states
     *   including session and contest states if applicable.
     * - Resets contest timers and fetches the updated total QSOs for the logbook in contest mode.
     * - Catches and handles any errors occurring during the logging process, ensuring appropriate
     *   debugging information is available.
     * - Resets the action state after completion.
     *
     * This function ensures the proper flow of the logging process, updating related states
     * and handling errors to provide a smooth user experience.
     *
     * @async
     * @function logContact
     * @returns {Promise<void>} A promise that resolves when the contact has been successfully logged.
     */
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

        updateQsoState();

        try {
            const qso: types.Qso = qsoState.toQso();
            qso.logbook_id = configState.logbook.id;

            await LogQso(qso);
            showToast.SUCCESS("QSO logged...");
            sessionState.update(await CurrentSessionQsoSlice());
            clipboardState.add(qso.comment);

            if ($isContestMode) {
                contestTimers.reset();
                contestState.totalQsos = await TotalQsosByLogbookId(configState.logbook.id);
            }
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