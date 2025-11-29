<script lang="ts">
    import {CALLSIGN_PATTERN} from "$lib/constants/callsign";
    import {handleAsyncError} from "$lib/utils/error-handler";
    import {NewQso} from "$lib/wailsjs/go/facade/Service";
    import {qsoState} from "$lib/states/qso-state.svelte";

    interface Props {
        id: string;
        label: string;
        value: string;
        labelCss?: string;
        divCss?: string;
        inputCss?: string;
        overallWidthCss?: string;
    }
    let {
        id,
        label,
        value = $bindable(),
        labelCss = 'block text-sm/5 font-medium',
        divCss = 'mt-2',
        inputCss = 'uppercase block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600',
        overallWidthCss = 'w-[150px]'
    }: Props = $props();

    let invalid = $state(false);
    let inputElement: HTMLInputElement;
    let lastKey: string | null = null;

    const isValid = (v: string): boolean => {
        const value = v.trim().toUpperCase();
        return CALLSIGN_PATTERN.test(value);
    }

    const handleInput = (e: Event): void => {
        const target = e.currentTarget as HTMLInputElement;
        if (!target) return;
        const v = target.value;
        if (v === '') {
            invalid = false;
            return;
        }
        invalid = !isValid(v);
    }

    const validateAndFocus = async (): Promise<void> => {
        const tabbed = lastKey === "Tab";
        lastKey = null;
        if (!tabbed) return;

        invalid = !isValid(value);
        if (invalid && inputElement) {
            inputElement.focus();
            inputElement.select();
            return;
        }

        try {
            const qso = await NewQso(value);
            console.log(">", qso);
            qsoState.createFromQSO(qso);
            qsoState.startTimer();
        } catch (e: unknown) {
            // Any error here is serious and means we cannot continue: either there is something wrong with the
            // provided callsign or the backend is not available.
            handleAsyncError(e, "Callsign.svelte: validateAndFocus");
            inputElement.focus();
            inputElement.select();
        } finally {

        }
    }

</script>

<div class={overallWidthCss}>
    <label for={id} class={labelCss}>{label}</label>
    <div class={divCss}>
        <input
                bind:this={inputElement}
                bind:value={value}
                type="text"
                id={id}
                class="{inputCss} {invalid ? 'outline-red-600 focus:outline-red-600' : ''}"
                autocomplete="off"
                spellcheck="false"
                aria-invalid={invalid}
                oninput={handleInput}
                onblur={validateAndFocus}
                onkeydown={(e) => lastKey = e.key}
        />
    </div>
</div>
