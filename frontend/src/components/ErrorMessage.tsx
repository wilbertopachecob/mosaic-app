import React from "react";
import { AlertCircle, AlertTriangle, Info, X } from "lucide-react";

interface ErrorMessageProps {
  message: string;
  onDismiss?: () => void;
  type?: "error" | "warning" | "info";
}

const iconMap = {
  error: AlertCircle,
  warning: AlertTriangle,
  info: Info,
};

const ErrorMessage: React.FC<ErrorMessageProps> = ({
  message,
  onDismiss,
  type = "error",
}) => {
  const Icon = iconMap[type];

  return (
    <div className={`alert alert-${type}`} role="alert">
      <div className="alert-content">
        <Icon size={18} aria-hidden="true" />
        <span>{message}</span>
      </div>
      {onDismiss && (
        <button
          type="button"
          className="icon-button"
          onClick={onDismiss}
          aria-label="Close"
        >
          <X size={16} aria-hidden="true" />
        </button>
      )}
    </div>
  );
};

export default ErrorMessage;
