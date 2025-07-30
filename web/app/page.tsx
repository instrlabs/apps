"use client";

import { useState } from "react";
import Link from "next/link";
import { ArrowRightIcon } from "@heroicons/react/24/outline";
import OutlinedButton from "@/components/outlined-button";
import { sidebar_menus } from "@/constants/routes";

export default function Home() {
  // Filter out only the image processing features (excluding Home, Histories, and dividers)
  const imageFeatures = sidebar_menus.filter(
    item => item.type === "text" && item.value !== "Home" && item.value !== "Histories"
  );

  return (
    <main className="pl-[300px] pt-[56px] min-h-screen">
      {/* Hero Section */}
      <section className="py-16 px-8 bg-gradient-to-r from-blue-50 to-indigo-50">
        <div className="max-w-4xl mx-auto">
          <h1 className="text-4xl font-bold mb-4">Powerful Image Processing Tools</h1>
          <p className="text-xl text-gray-700 mb-8">
            Transform your images with our easy-to-use tools. Compress, resize, crop, convert, and
            rotate your images in seconds.
          </p>
          <div className="flex space-x-4">
            <Link href="/compress-image">
              <OutlinedButton className="bg-blue-600 text-white border-blue-600 hover:bg-blue-700 hover:border-blue-700 px-6 py-2">
                Get Started
              </OutlinedButton>
            </Link>
            <Link href="/histories">
              <OutlinedButton>View History</OutlinedButton>
            </Link>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 px-8">
        <div className="max-w-4xl mx-auto">
          <h2 className="text-2xl font-bold mb-8 text-center">Our Image Processing Features</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {imageFeatures.map((feature, index) => (
              <Link key={index} href={feature.type === "text" ? feature.href : "#"}>
                <div className="border border-gray-200 rounded-lg p-6 hover:shadow-md transition-shadow">
                  <h3 className="text-xl font-semibold mb-2">{feature.value}</h3>
                  <p className="text-gray-600 mb-4">{getFeatureDescription(feature.value)}</p>
                  <div className="flex items-center text-blue-600">
                    <span className="mr-2">Try now</span>
                    <ArrowRightIcon className="h-4 w-4" />
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* Call to Action */}
      <section className="py-12 px-8 bg-gray-100">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-2xl font-bold mb-4">Ready to transform your images?</h2>
          <p className="text-gray-700 mb-6">
            Our tools are free, fast, and easy to use. No registration required.
          </p>
          <Link href="/compress-image">
            <OutlinedButton className="bg-blue-600 text-white border-blue-600 hover:bg-blue-700 hover:border-blue-700 px-6 py-2">
              Start Processing
            </OutlinedButton>
          </Link>
        </div>
      </section>
    </main>
  );
}

// Helper function to get feature descriptions
function getFeatureDescription(featureName: any) {
  switch (featureName) {
    case "Compress":
      return "Reduce file size while maintaining quality";
    case "Resize":
      return "Change dimensions of your images";
    case "Crop":
      return "Remove unwanted areas from your images";
    case "Convert":
      return "Change image format (JPG, PNG, WebP, etc.)";
    case "Rotate":
      return "Rotate or flip your images";
    default:
      return "Process your images with ease";
  }
}
