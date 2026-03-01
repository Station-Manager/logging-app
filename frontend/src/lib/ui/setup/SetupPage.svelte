<script lang="ts">
    import {inputBase, inputBaseUppercase, labelBase} from "@station-manager/shared-utils";

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

</script>

<div class="mx-6 flex flex-col gap-4">
    <div>
        <h2 class="text-center">Setting your Default Logbook</h2>
        <p>This page will create the default logbook used by the Station Manager logging application. The callsign
            you provided here will identify the <b>Station Callsign</b> for all QSOs associated with this log book.
        </p>
        <p>If you use QRZ.com, then the callsign entered here should be the same as the callsign for the QRZ.com log book
            to which the QSOs will be forwarded (forwarding of QSOs is configurable and not enabled by default).
        </p>
    </div>
    <div>
        <div class="flex flex-col w-37.5">
            <label for="name">Name</label>
            <input
                    bind:value={logbook.name}
                    id="name"
                    type="text"
                    class={inputBase}
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
        <div class="flex flex-col w-48">
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
</div>
