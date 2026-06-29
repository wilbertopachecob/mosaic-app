import React from "react";
import { useTranslation } from "react-i18next";
import { AlertCircle, AlertTriangle, Info, X } from "lucide-react";

/** Props for a dismissible inline alert banner. */
type ErrorMessageProps = {
  message: string;
  onDismiss?: () => void;
  type?: "error" | "warning" | "info";
};

const iconMap = {
  error: AlertCircle,
  warning: AlertTriangle,
  info: Info,
};

/**
 * Inline alert for errors, warnings, or informational messages.
 */
const ErrorMessage: React.FC<ErrorMessageProps> = ({
  message,
  onDismiss,
  type = "error",
}) => {
  const { t } = useTranslation();
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
          aria-label={t("error.close")}
        >
          <X size={16} aria-hidden="true" />
        </button>
      )}
    </div>
  );
};

export default ErrorMessage;
