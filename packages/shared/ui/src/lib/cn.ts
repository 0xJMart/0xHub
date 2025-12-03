export type ClassValue = string | number | null | undefined | ClassValue[];

export const cn = (...values: ClassValue[]): string => {
  const classes: string[] = [];

  const visit = (value: ClassValue): void => {
    if (value === null || value === undefined) {
      return;
    }
    if (Array.isArray(value)) {
      value.forEach(visit);
      return;
    }
    if (typeof value === 'string' || typeof value === 'number') {
      classes.push(String(value));
      return;
    }
  };

  values.forEach(visit);

  return classes.filter(Boolean).join(' ');
};


