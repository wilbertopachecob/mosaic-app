import React from 'react';

/** Size variants for the Bootstrap spinner. */
type LoadingSpinnerProps = {
  size?: 'sm' | 'md' | 'lg';
  color?: string;
};

/**
 * Bootstrap-styled loading spinner.
 * @deprecated Not currently used in the app; prefer inline `Loader2` spinners.
 */
const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({ 
  size = 'md', 
  color = 'primary' 
}) => {
  const sizeClasses = {
    sm: 'spinner-border-sm',
    md: '',
    lg: 'spinner-border-lg'
  };

  return (
    <div className={`spinner-border text-${color} ${sizeClasses[size]}`} role="status">
      <span className="visually-hidden">Loading...</span>
    </div>
  );
};

export default LoadingSpinner;
