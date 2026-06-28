import React from "react";
import { useTranslation } from "react-i18next";
import { useTheme } from "../context/ThemeContext";
import logoLockupDark from "../assets/logo-lockup-dark.svg";
import logoLockupLight from "../assets/logo-lockup-light.svg";

const Logo: React.FC = () => {
  const { theme } = useTheme();
  const { t } = useTranslation();
  const logoSrc = theme === "dark" ? logoLockupDark : logoLockupLight;

  return (
    <h1 className="app-logo">
      <img src={logoSrc} alt={t("header.title")} className="app-logo-image" />
    </h1>
  );
};

export default Logo;
