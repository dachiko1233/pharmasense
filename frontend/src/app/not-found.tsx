import Link from "next/link";

export default function RootNotFound() {
  return (
    <div
      style={{
        display: "flex",
        minHeight: "100vh",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        backgroundColor: "#f8fafc",
        fontFamily: "sans-serif",
      }}
    >
      <div style={{ textAlign: "center" }}>
        <h1 style={{ fontSize: "4rem", fontWeight: "bold", color: "#0f172a", margin: 0 }}>
          404
        </h1>
        <p style={{ marginTop: "1rem", fontSize: "1.25rem", color: "#475569" }}>
          Page not found
        </p>
        <Link
          href="/en"
          style={{
            display: "inline-block",
            marginTop: "1.5rem",
            padding: "0.5rem 1rem",
            backgroundColor: "#2563eb",
            color: "#fff",
            borderRadius: "0.375rem",
            textDecoration: "none",
          }}
        >
          Go home
        </Link>
      </div>
    </div>
  );
}
