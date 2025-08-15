import React from "react";

interface FormInputProps {
  id: string;
  type: string;
  label: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  placeholder: string;
}

const FormInput: React.FC<FormInputProps> = ({
  id,
  type,
  value,
  onChange,
  placeholder,
}) => (
  <input
    id={id}
    type={type}
    placeholder={placeholder}
    className="w-full px-5 py-4 bg-white shadow-primary rounded-full text-sm outline-none"
    value={value}
    onChange={onChange}
    required
  />
);

export default FormInput;
