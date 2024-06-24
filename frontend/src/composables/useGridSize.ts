import {onBeforeUnmount, onMounted, ref} from 'vue';
import {StoryModule} from "../../bindings/storyguardian/src/project";

const breakpoints = {
    sm: '(min-width: 640px)',
    md: '(min-width: 768px)',
    lg: '(min-width: 1024px)',
    xl: '(min-width: 1280px)',
};

function getCurrentBreakpoint(): string {
    const sortedBreakpoints = Object.entries(breakpoints).sort(
        ([, a], [, b]) => parseInt(a.match(/\d+/)![0]) - parseInt(b.match(/\d+/)![0])
    );

    for (const [breakpoint, mediaQuery] of sortedBreakpoints.reverse()) {
        if (window.matchMedia(mediaQuery).matches) {
            return breakpoint;
        }
    }

    return 'sm';
}

export function useGridSize(moduleConfig: StoryModule | undefined) {
    const columnSize = ref('col-span-1');

    function updateColumnSize(newColumnSize: string) {
        const currentBreakpoint = getCurrentBreakpoint();
        const maxColumns = {
            sm: 1,
            md: 2,
            lg: 3,
            xl: 4,
        }[currentBreakpoint];

        if(!maxColumns){
            columnSize.value = `col-span-1`;
            return;
        }

        const adjustedColumnSize = Math.min(Number(newColumnSize), maxColumns).toString();
        columnSize.value = `col-span-${adjustedColumnSize}`;
    }

    function changeGridSize(module: string, newColumnSize: string, emit: any) {
        if (!moduleConfig) return;

        moduleConfig.configuration['columnSize'] = newColumnSize;
        emit('configChange', module, 'columnSize', newColumnSize);

        updateColumnSize(newColumnSize);
    }

    function handleResize() {
        if (!moduleConfig) return;
        updateColumnSize(moduleConfig.configuration['columnSize']);
    }

    onMounted(() => {
        if (!moduleConfig) return;
        columnSize.value = `col-span-${moduleConfig.configuration['columnSize']}`;

        Object.values(breakpoints).forEach((mediaQuery) => {
            const mql = window.matchMedia(mediaQuery);
            mql.addEventListener('change', handleResize);
        });
    });

    onBeforeUnmount(() => {
        Object.values(breakpoints).forEach((mediaQuery) => {
            const mql = window.matchMedia(mediaQuery);
            mql.removeEventListener('change', handleResize);
        });
    });

    return {
        columnSize,
        changeGridSize,
    };
}