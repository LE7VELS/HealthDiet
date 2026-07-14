import { Badge as AstryxBadge, type BadgeProps } from '@astryxdesign/core/Badge'
import { AppShell as AstryxAppShell, type AppShellProps } from '@astryxdesign/core/AppShell'
import {
  CheckboxInput as AstryxCheckboxInput,
  type CheckboxInputProps,
} from '@astryxdesign/core/CheckboxInput'
import {
  DropdownMenu as AstryxDropdownMenu,
  DropdownMenuItem as AstryxDropdownMenuItem,
  type DropdownMenuProps,
  type DropdownMenuItemProps,
} from '@astryxdesign/core/DropdownMenu'
import { Heading as AstryxHeading, type HeadingProps } from '@astryxdesign/core/Heading'
import { Grid as AstryxGrid, type GridProps } from '@astryxdesign/core/Grid'
import { HStack as AstryxHStack, type HStackProps } from '@astryxdesign/core/HStack'
import { Icon as AstryxIcon, type IconProps } from '@astryxdesign/core/Icon'
import {
  ProgressBar as AstryxProgressBar,
  type ProgressBarProps,
} from '@astryxdesign/core/ProgressBar'
import {
  SideNav as AstryxSideNav,
  SideNavItem as AstryxSideNavItem,
  type SideNavProps,
  type SideNavItemProps,
} from '@astryxdesign/core/SideNav'
import { Text as AstryxText, type TextProps } from '@astryxdesign/core/Text'
import { VStack as AstryxVStack, type VStackProps } from '@astryxdesign/core/VStack'

// 业务代码只依赖项目级封装，后续升级 Astryx 时可在这里集中适配 API 变化。
export function Badge(props: BadgeProps) {
  return <AstryxBadge {...props} />
}

export function ApplicationShell(props: AppShellProps) {
  return <AstryxAppShell {...props} />
}

export function CheckboxInput(props: CheckboxInputProps) {
  return <AstryxCheckboxInput {...props} />
}

export function DropdownMenu(props: DropdownMenuProps) {
  return <AstryxDropdownMenu {...props} />
}

export function DropdownMenuItem(props: DropdownMenuItemProps) {
  return <AstryxDropdownMenuItem {...props} />
}

export function Heading(props: HeadingProps) {
  return <AstryxHeading {...props} />
}

export function Grid(props: GridProps) {
  return <AstryxGrid {...props} />
}

export function HStack(props: HStackProps) {
  return <AstryxHStack {...props} />
}

export function Icon(props: IconProps) {
  return <AstryxIcon {...props} />
}

export function ProgressBar(props: ProgressBarProps) {
  return <AstryxProgressBar {...props} />
}

export function SideNav(props: SideNavProps) {
  return <AstryxSideNav {...props} />
}

export function SideNavItem(props: SideNavItemProps) {
  return <AstryxSideNavItem {...props} />
}

export function Text(props: TextProps) {
  return <AstryxText {...props} />
}

export function VStack(props: VStackProps) {
  return <AstryxVStack {...props} />
}
