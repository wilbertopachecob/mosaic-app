import React from 'react';

interface ErrorMessageProps {
  message: string;
  onDismiss?: () => void;
  type?: 'error' | 'warning' | 'info';
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ 
  message, 
  onDismiss, 
  type = 'error' 
}) => {
  const alertClasses = {
    error: 'alert-danger',
    warning: 'alert-warning',
    info: 'alert-info'
  };

  const iconClasses = {
    error: 'fas fa-exclamation-circle',
    warning: 'fas fa-exclamation-triangle',
    info: 'fas fa-info-circle'
  };

  return (
    <div className={`alert ${alertClasses[type]} alert-dismissible fade show`} role="alert">
      <div className="d-flex align-items-center">
        <i className={`${iconClasses[type]} me-2`}></i>
        <span>{message}</span>
      </div>
      {onDismiss && (
        <button
          type="button"
          className="btn-close"
          onClick={onDismiss}
          aria-label="Close"
        ></button>
      )}
    </div>
  );
};

export default ErrorMessage; 