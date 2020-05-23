import React from 'react';

import emotionStyled, { CreateStyled } from '@emotion/styled/macro';
import { ThemeProvider as EmotionThemeProvider } from 'emotion-theming';

export interface Theme {
  fonts: {
    light: string;
    regular: string;
    semiBold: string;
    bold: string;
    extraBold: string;
  };
  sizes: {
    s: string;
    m: string;
    l: string;
    xl: string;
  };
  colors: {
    blue: string;
    darkBlue: string;
    gray: string;
    darkGray: string;
    white: string;
    whitish: string;
    pink: string;
    green: string;
    orange: string;
    tileBack: string;
    purple: string;
    lightPurple: string;
  };
}

const theme: Theme = {
  fonts: {
    light: "'OpenSans Light'",
    regular: "'OpenSans Regular'",
    semiBold: "'OpenSans SemiBold'",
    bold: "'OpenSans Bold'",
    extraBold: "'OpenSans ExtraBold'",
  },
  sizes: {
    s: '14px',
    m: '18px',
    l: '22px',
    xl: '27px',
  },
  colors: {
    blue: '#252f4a',
    darkBlue: '#212133',
    gray: '#848a99',
    darkGray: '#6b6969',
    white: '#ffffff',
    whitish: '#f5f5f5',
    pink: '#f5406e',
    green: '#46E80E',
    orange: '#f66b1c',
    tileBack: 'rgba(245,245,245,0.04)',
    purple: '#57038d',
    lightPurple: '#a540cd',
  },
};

export const styled = emotionStyled as CreateStyled<Theme>;

export const ThemeProvider: React.FC = ({ children }) => {
  return <EmotionThemeProvider theme={theme}>{children}</EmotionThemeProvider>;
};