import type { LayoutData } from './$types';

export const prerender = true
export const ssr = false

export const load: LayoutData = async (): Promise<object> => {
    return {}
}
