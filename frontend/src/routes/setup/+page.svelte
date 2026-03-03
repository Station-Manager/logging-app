<script lang="ts">
    import {inputBase, inputBaseUppercase, labelBase} from "@station-manager/shared-utils";
    import {types} from "$lib/wailsjs/go/models";
    import {FinaliseSetup} from "$lib/wailsjs/go/facade/Service";

    interface Logbook {
        name: string;
        callsign: string;
        description: string;
    }

    let logbook: Logbook = $state({
        name: "Default",
        callsign: "",
        description: "",
    });

    const onClickSave = async (): Promise<void> => {
        try {
            const lb: types.Logbook = new types.Logbook();
            lb.name = logbook.name;
            lb.callsign = logbook.callsign;
            lb.description = logbook.description;
            await FinaliseSetup(lb);
        } catch(e: unknown) {
            console.error(e);
        }
    }
</script>

<div class="mx-6 flex flex-col gap-4">
    <div>
        <h2 class="text-center">Setting your Default Logbook</h2>
        <p>This page will create the <b>default logbook</b> used by the <i>Station Manager</i> logging application. The callsign
            you provide here will identify the <i>Station Callsign</i> for all QSOs associated with this log book.
        </p>
        <p>If you use QRZ.com, then the callsign entered here should be the same as the callsign for the QRZ.com log book
            to which the QSOs will be forwarded (forwarding of QSOs is configurable and not enabled by default).
        </p>
    </div>
    <div class="flex flex-col gap-y-4 mx-8">
        <div>
            <div class="flex flex-col w-37.5">
                <label for="name" class={labelBase}>Name</label>
                <input
                        bind:value={logbook.name}
                        id="name"
                        type="text"
                        class="{inputBase} text-gray-400"
                        title="This cannot be changed."
                        disabled>
            </div>
        </div>
        <div>
            <div class="flex flex-col w-37.5">
                <label class={labelBase} for="callsign">Callsign</label>
                <input
                        bind:value={logbook.callsign}
                        id="callsign"
                        type="text"
                        placeholder="Callsign"
                        class={inputBaseUppercase}
                />
            </div>
        </div>
        <div>
            <div class="flex flex-col w-64">
                <label class={labelBase} for="description">Description</label>
                <input
                        bind:value={logbook.description}
                        id="description"
                        type="text"
                        placeholder="Description"
                        class={inputBase}
                />
            </div>
        </div>
        <div class="flex justify-end">
            <button
                    onclick={onClickSave}
                    class="h-9 cursor-pointer rounded-md bg-indigo-600 p-2.5 py-1.5 text-base font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 w-20">Save</button>
        </div>
    </div>
</div>
