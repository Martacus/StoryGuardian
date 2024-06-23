import {onMounted, ref} from 'vue';
import {StoryModule} from "../../bindings/storyguardian/internal/project";

export function useGridSize(moduleConfig: StoryModule | undefined) {
    const columnSize = ref('col-span-4');

    function changeGridSize(newColumnSize: string, emit: any) {
        if (!moduleConfig) return;
        moduleConfig.configuration['columnSize'] = newColumnSize;
        columnSize.value = columnSize.value.slice(0, -1) + newColumnSize;
        emit('configChange', 'columnSize', newColumnSize);
    }

    onMounted(() => {
        if (!moduleConfig) return;
        columnSize.value = columnSize.value.slice(0, -1) + moduleConfig.configuration['columnSize'];
    });

    return {
        columnSize,
        changeGridSize,
    };
}