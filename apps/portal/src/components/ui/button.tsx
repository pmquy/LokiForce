import * as React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '../../lib/utils';

const buttonVariants = cva(
  'inline-flex items-center justify-center whitespace-nowrap rounded-2xl text-sm font-bold transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500/20 disabled:pointer-events-none disabled:opacity-50 active:scale-98 cursor-pointer select-none',
  {
    variants: {
      variant: {
        default: 'bg-gradient-to-r from-indigo-500 to-teal-500 text-white shadow-lg shadow-indigo-500/20 hover:opacity-95',
        destructive: 'bg-rose-600 text-white hover:bg-rose-700 shadow-lg shadow-rose-600/20',
        outline: 'border border-slate-800 bg-slate-950 text-slate-200 hover:bg-slate-900 hover:text-white',
        secondary: 'bg-slate-800 text-slate-100 hover:bg-slate-700',
        ghost: 'text-slate-400 hover:bg-slate-900 hover:text-white',
        link: 'text-indigo-400 underline-offset-4 hover:underline',
      },
      size: {
        default: 'h-11 px-5 py-2.5',
        sm: 'h-9 rounded-xl px-3 text-xs',
        lg: 'h-12 rounded-2xl px-6',
        icon: 'h-11 w-11',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : 'button';
    return (
      <Comp
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      />
    );
  }
);
Button.displayName = 'Button';

export { Button, buttonVariants };
