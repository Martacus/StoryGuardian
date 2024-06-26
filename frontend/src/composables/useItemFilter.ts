import {Ref, ref, watch} from 'vue';

export function useItemFilter<T>(items: Ref<T[]>, filterFn: (item: T, filter: string) => boolean) {
    const searchInput = ref('');
    const searchResult = ref<T[]>([]) as Ref<T[]>;


    watch(searchInput, async (newInput) => {
        searchResult.value = items.value.filter((item) => filterFn(item, newInput));
    });

    return {
        searchInput,
        searchResult,
    };
}