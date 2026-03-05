import { defineAsyncComponent, type Component } from 'vue'
import { FileText, Users, Image as ImageIcon } from 'lucide-vue-next'
import { ModuleLayout } from '../../bindings/litguardian/models'

export interface ModuleDefinition {
  id: string
  label: string
  icon: Component
  defaultColSpan: number
  component: Component
  views: string[]
}

export const moduleRegistry: ModuleDefinition[] = [
  {
    id: 'overview-info',
    label: 'World Info',
    icon: FileText,
    defaultColSpan: 6,
    component: defineAsyncComponent(() => import('./overview/InfoModule.vue')),
    views: ['overview'],
  },
  {
    id: 'overview-entities',
    label: 'Entities',
    icon: Users,
    defaultColSpan: 6,
    component: defineAsyncComponent(() => import('./overview/EntitiesModule.vue')),
    views: ['overview'],
  },
  {
    id: 'overview-images',
    label: 'Images',
    icon: ImageIcon,
    defaultColSpan: 12,
    component: defineAsyncComponent(() => import('./overview/ImagesModule.vue')),
    views: ['overview'],
  },
]

/**
 * Returns the default ModuleLayout array for a view when no layout.json exists.
 * Filters the registry by viewId and maps each definition to its defaults.
 */
export function getDefaultModules(viewId: string): ModuleLayout[] {
  return moduleRegistry
    .filter(m => m.views.includes(viewId))
    .map(m => new ModuleLayout({ id: m.id, colSpan: m.defaultColSpan, visible: true }))
}
