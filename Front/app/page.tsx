import CornPurchase from "@/components/corn-purchase"
import { CropIcon as Corn } from "lucide-react"

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-4 bg-amber-50">
      <div className="w-full max-w-md p-6 bg-white rounded-lg shadow-md border border-amber-200">
        <div className="flex items-center justify-center mb-6">
          <Corn className="h-10 w-10 text-yellow-500 mr-2" />
          <h1 className="text-2xl font-bold text-amber-800">Bob's Corn Shop</h1>
        </div>

        <p className="text-center mb-6 text-amber-700">Get the freshest corn straight from Bob's farm!</p>

        <CornPurchase />
      </div>
    </main>
  )
}
