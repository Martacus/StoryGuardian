import {onMounted, Ref, ref, watch} from "vue";

export function useBasicListHeight(itemView: Ref<string, string>, searchResult: Ref<string[]>) {
    const listHeight = ref<string>('h-0');

    function calcListHeight() {
        if (itemView.value === 'list') {
            if (searchResult.value.length > 8) {
                listHeight.value = 'h-96';
            } else {
                listHeight.value = 'h-' + Math.max(searchResult.value.length, 1) * 12;
            }
        } else {
            if (searchResult.value.length > 24) {
                listHeight.value = 'h-96';
            } else {
                listHeight.value = 'h-' + Math.max(searchResult.value.length, 3) / 3 * 12;
            }
        }
    }

    watch([itemView, searchResult], () => {
        calcListHeight();
    });

    onMounted(() => {
    });

    return {
        listHeight
    };
}