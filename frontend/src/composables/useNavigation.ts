import { useRouter } from 'vue-router';

export function useNavigation() {
    const router = useRouter();

    async function navigateToEntity(id: string) {
        await router.push('/entity/' + id);
    }

    return {
        navigateToEntity,
    };
}