import {Ref, ref, watch} from 'vue';

export function useItemFilter<T>(items: Ref<T[]>, filterFn: (item: T, filter: string) => boolean) {
    const searchInput = ref('');
    const searchResult = ref<T[]>([]) as Ref<T[]>;

    watch(searchInput, async (newInput) => {
        searchResult.value = items.value.filter((item) => filterFn(item, newInput));
    });
    watch(items, async (newItems) => {
        searchResult.value = newItems.filter((item) => filterFn(item, searchInput.value));
    });

    return {
        searchInput,
        searchResult,
    };
}