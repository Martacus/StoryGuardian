import {onMounted, ref} from 'vue';
import {StoryModule} from "../../bindings/storyguardian/src/project";

type ListViewMode = 'grid' | 'list';
export function useItemGridLayout(moduleConfig: StoryModule | undefined) {
    const itemView = ref<ListViewMode>('list');
    function changeItemView(module: string, newItemView: string, emit: any) {
        if (!moduleConfig) return;

        if(itemView.value === 'list'){
            itemView.value = 'grid';
        } else {
            itemView.value = 'list';
        }
        emit('configChange', module, 'itemView', newItemView);
    }


    onMounted(() => {
        if (!moduleConfig) return;
        itemView.value = moduleConfig.configuration['itemView'] as ListViewMode;
    });

    return {
        itemView,
        changeItemView,
    };
}