import { defineAsyncComponent, type Component } from 'vue'
import { FileText, Users, Image as ImageIcon, Info, Link } from 'lucide-vue-next'
import { ModuleLayout } from '../../bindings/litguardian/models'

export interface ModuleDefinition {
  id: string
  label: string
  icon: Component
  defaultW: number
  defaultH: number
  minW?: number
  minH?: number
  component: Component
  views: string[]
}

export const moduleRegistry: ModuleDefinition[] = [
  {
    id: 'overview-info',
    label: 'World Info',
    icon: FileText,
    defaultW: 6,
    defaultH: 8,
    minW: 3,
    minH: 3,
    component: defineAsyncComponent(() => import('./overview/InfoModule.vue')),
    views: ['overview'],
  },
  {
    id: 'overview-entities',
    label: 'Entities',
    icon: Users,
    defaultW: 6,
    defaultH: 6,
    minW: 3,
    minH: 3,
    component: defineAsyncComponent(() => import('./overview/EntitiesModule.vue')),
    views: ['overview'],
  },
  {
    id: 'overview-images',
    label: 'Images',
    icon: ImageIcon,
    defaultW: 12,
    defaultH: 6,
    minW: 3,
    minH: 3,
    component: defineAsyncComponent(() => import('./overview/ImagesModule.vue')),
    views: ['overview'],
  },
  {
    id: 'entity-info',
    label: 'Basic Info',
    icon: Info,
    defaultW: 6,
    defaultH: 10,
    minW: 3,
    minH: 4,
    component: defineAsyncComponent(() => import('./entity/InfoModule.vue')),
    views: ['entity-detail'],
  },
  {
    id: 'entity-relations',
    label: 'Relations',
    icon: Link,
    defaultW: 6,
    defaultH: 10,
    minW: 3,
    minH: 4,
    component: defineAsyncComponent(() => import('./entity/RelationsModule.vue')),
    views: ['entity-detail'],
  },
]

/**
 * Returns the default ModuleLayout array for a view when no layout.json exists.
 * Auto-positions modules left-to-right, wrapping to the next row when x + w > 12.
 */
export function getDefaultModules(viewId: string): ModuleLayout[] {
  const defs = moduleRegistry.filter(m => m.views.includes(viewId))
  const result: ModuleLayout[] = []
  let x = 0
  let y = 0
  let rowMaxH = 0

  for (const def of defs) {
    if (x + def.defaultW > 12) {
      x = 0
      y += rowMaxH
      rowMaxH = 0
    }
    result.push(new ModuleLayout({
      id: def.id,
      x,
      y,
      w: def.defaultW,
      h: def.defaultH,
      visible: true,
    }))
    rowMaxH = Math.max(rowMaxH, def.defaultH)
    x += def.defaultW
  }

  return result
}
