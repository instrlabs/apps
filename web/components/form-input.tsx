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
  label,
  value,
  onChange,
  placeholder,
}) => (
  <div className="space-y-1">
    <label htmlFor={id} className="text-sm font-medium">
      {label}
    </label>
    <input
      id={id}
      type={type}
      placeholder={placeholder}
      className="px-2 py-2.5 rounded w-full outline-none text-sm border border-gray-300"
      value={value}
      onChange={onChange}
      required
    />
  </div>
);

export default FormInput;