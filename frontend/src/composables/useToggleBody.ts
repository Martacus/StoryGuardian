import {onMounted, ref} from 'vue';
import {StoryModule} from "../../bindings/storyguardian/src/project";

export function useToggleBody(moduleConfig: StoryModule | undefined) {
    const showCardBody = ref(true);

    function toggleCardBody(module: string, emit: any) {
        showCardBody.value = !showCardBody.value;
        emit('configChange', module, 'open', JSON.stringify(showCardBody.value))
    }

    onMounted(() => {
        if (!moduleConfig) return;
        showCardBody.value = moduleConfig.configuration['open'] === 'true';
    });

    return {
        showCardBody,
        toggleCardBody,
    };
}