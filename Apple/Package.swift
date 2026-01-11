// swift-tools-version: 5.10
import PackageDescription

let package = Package(
    name: "ExqcorCore",
    platforms: [
        .macOS(.v14),
        .iOS(.v17)
    ],
    products: [
        .library(
            name: "ExqcorCore",
            targets: ["ExqcorCore"]
        )
    ],
    targets: [
        .target(
            name: "ExqcorCore",
            path: "Sources/Core"
        )
    ]
)
