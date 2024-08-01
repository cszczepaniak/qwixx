<script lang="ts">
	import clsx from "clsx";
	import { fade } from "svelte/transition";

	type Num = "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" | "10" | "11" | "12" | "LOCK";

	let allNums: Num[] = ["2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"];
	let numsAsc = [...allNums, "LOCK"];
	let numsDesc = [...[...allNums].reverse(), "LOCK"];

	let rows = [numsAsc, numsAsc, numsDesc, numsDesc];

	let colors = [
		{ bg: "bg-red-100", fg: "bg-red-800", mg: "bg-red-500", text: "text-red-500" },
		{ bg: "bg-yellow-100", fg: "bg-yellow-600", mg: "bg-yellow-300", text: "text-yellow-300" },
		{ bg: "bg-green-100", fg: "bg-green-700", mg: "bg-green-500", text: "text-green-500" },
		{ bg: "bg-blue-100", fg: "bg-blue-900", mg: "bg-blue-700", text: "text-blue-700" },
	];

	let confirmClearCrossDialog: HTMLDialogElement;
	let confirmClearAllDialog: HTMLDialogElement;
	let rowColToConfirm = $state([-1, -1]);

	let scoreDialog: HTMLDialogElement;

	let crosses = $state([
		[false, false, false, false, false, false, false, false, false, false, false, false],
		[false, false, false, false, false, false, false, false, false, false, false, false],
		[false, false, false, false, false, false, false, false, false, false, false, false],
		[false, false, false, false, false, false, false, false, false, false, false, false],
	]);

	function countCrossesInRow(row: boolean[]) {
		return row.reduce((ct, v) => ct + (v ? 1 : 0), 0);
	}

	let missedRolls = $state([false, false, false, false]);

	let score = $derived(
		(() => {
			const scoresByXCount = [0, 1, 3, 6, 10, 15, 21, 28, 36, 45, 55, 66, 78];

			let countPerRow = crosses.map(countCrossesInRow);
			let scoreForCrosses = countPerRow.reduce((ct, v) => scoresByXCount[v] + ct, 0);
			console.log(scoreForCrosses);
			let nMissedRolls = countCrossesInRow(missedRolls);

			return scoreForCrosses - 5 * nMissedRolls;
		})(),
	);

	let unlockedLocks = $derived(crosses.map((c) => countCrossesInRow(c) >= 5));
	let disabledIndices = $derived(crosses.map((c) => c.lastIndexOf(true)));
</script>

<div class="flex h-screen flex-row items-center justify-around bg-gray-200 shadow-sm">
	<div class="flex flex-col space-y-0.5">
		{#each rows as row, rowIdx}
			<div class="relative">
				{#if !unlockedLocks[rowIdx]}
					<div
						class="something translate diagonal-stripe-1 absolute z-10 h-[3.75rem]
					w-[7rem] translate-x-[32.25rem] translate-y-[0.25rem] rounded-md border-2 border-dashed border-gray-800"
						transition:fade={{ delay: 250, duration: 300 }}
					></div>
				{/if}
				<div class={clsx("mx-auto w-max px-3 py-2", colors[rowIdx].mg)}>
					<div
						class={clsx(
							"flex flex-row items-center space-x-0.5 rounded-md p-0.5",
							colors[rowIdx].fg,
						)}
					>
						{#each row as col, colIdx}
							<div class="relative">
								{#if crosses[rowIdx][colIdx]}
									<button
										class="absolute left-1/2 top-1/2 h-12 w-12 -translate-x-1/2 -translate-y-1/2 transform"
										onclick={() => {
											rowColToConfirm = [rowIdx, colIdx];
											confirmClearCrossDialog.showModal();
										}}
									>
										{@render xIcon("size-12")}
									</button>
								{/if}
								<button
									onclick={() => (crosses[rowIdx][colIdx] = true)}
									disabled={colIdx <= disabledIndices[rowIdx]}
									class={clsx(
										"flex h-12 w-12 items-center justify-around rounded-md p-2 text-center text-2xl font-bold",
										colors[rowIdx].bg,
										colors[rowIdx].text,
										// The 10th index is either 2 or 12, and it's where we
										// separate the lock from the rest of the squares.
										colIdx === 10 && "ml-2",
									)}
								>
									{#if col === "LOCK"}
										{@render lockIcon()}
									{:else}
										{col}
									{/if}
								</button>
							</div>
						{/each}
					</div>
				</div>
			</div>
		{/each}

		<div>
			<div class="mt-2 flex flex-row justify-between">
				<div class="flex flex-row space-x-2">
					<button
						onclick={() => scoreDialog.showModal()}
						class="rounded-md border border-gray-800 px-4 py-2"
					>
						Show Score
					</button>
					<button
						onclick={() => confirmClearAllDialog.showModal()}
						class="rounded-md border border-gray-800 px-4 py-2"
					>
						Clear Scorecard
					</button>
				</div>
				<div class="flex flex-col items-end space-y-2">
					<p class="font-medium italic">Missed Rolls (-5 pts)</p>
					<div class="flex flex-row space-x-2">
						{#each missedRolls as missed, i}
							<div class="relative">
								{#if missed}
									<button
										class="absolute left-1/2 top-1/2 h-6 w-6 -translate-x-1/2 -translate-y-1/2 transform"
										onclick={() => (missedRolls[i] = false)}
									>
										{@render xIcon("size-6")}
									</button>
								{/if}
								<button
									onclick={() => (missedRolls[i] = true)}
									class="flex h-6 w-6 items-center justify-around rounded-md border
							border-gray-800 bg-gray-100 p-2 text-center text-2xl font-bold"
								></button>
							</div>
						{/each}
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<dialog bind:this={confirmClearCrossDialog} class="space-y-6 rounded-md p-6">
	<p class="font-medium">Really clear this cross?</p>
	<div class="flex flex-row justify-between space-x-2">
		<button
			class="rounded-md border px-4 py-2"
			onclick={() => {
				let [row, col] = rowColToConfirm;
				crosses[row][col] = false;
				confirmClearCrossDialog.close();
			}}
		>
			Confirm
		</button>
		<button class="rounded-md border px-4 py-2" onclick={() => confirmClearCrossDialog.close()}>
			Cancel
		</button>
	</div>
</dialog>

<dialog bind:this={confirmClearAllDialog} class="space-y-6 rounded-md p-6">
	<p class="font-medium">Really clear the entire scorecard?</p>
	<div class="flex flex-row justify-between space-x-2">
		<button
			class="rounded-md border px-4 py-2"
			onclick={() => {
				for (let i = 0; i < crosses.length; i++) {
					for (let j = 0; j < crosses[i].length; j++) {
						crosses[i][j] = false;
					}
				}
				for (let i = 0; i < missedRolls.length; i++) {
					missedRolls[i] = false;
				}
				confirmClearAllDialog.close();
			}}
		>
			Confirm
		</button>
		<button class="rounded-md border px-4 py-2" onclick={() => confirmClearAllDialog.close()}>
			Cancel
		</button>
	</div>
</dialog>

<dialog bind:this={scoreDialog} class="rounded-md">
	<div class="flex flex-row justify-end pr-2 pt-2">
		<button onclick={() => scoreDialog.close()}>{@render xIcon("size-6")}</button>
	</div>
	<p class="p-6 text-2xl">Your score is: {score}</p>
</dialog>

{#snippet xIcon(size)}
	<svg
		xmlns="http://www.w3.org/2000/svg"
		fill="none"
		viewBox="0 0 24 24"
		stroke-width="1"
		stroke="currentColor"
		class={size ?? "size-6"}
	>
		<path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
	</svg>
{/snippet}

{#snippet lockIcon()}
	<svg
		xmlns="http://www.w3.org/2000/svg"
		fill="none"
		viewBox="0 0 24 24"
		stroke-width="1.5"
		stroke="currentColor"
		class="size-6"
	>
		<path
			stroke-linecap="round"
			stroke-linejoin="round"
			d="M13.5 10.5V6.75a4.5 4.5 0 1 1 9 0v3.75M3.75 21.75h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H3.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25Z"
		/>
	</svg>
{/snippet}
