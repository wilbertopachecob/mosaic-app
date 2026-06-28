import { screen, fireEvent } from "@testing-library/react";
import "@testing-library/jest-dom";
import ErrorMessage from "@/components/ErrorMessage";
import { renderWithProviders } from "@/testUtils";

describe("ErrorMessage", () => {
  it("renders the message with alert role", () => {
    renderWithProviders(<ErrorMessage message="Something went wrong" />);

    const alert = screen.getByRole("alert");
    expect(alert).toHaveTextContent("Something went wrong");
    expect(alert).toHaveClass("alert-error");
  });

  it("applies warning and info variants", () => {
    const { rerender } = renderWithProviders(
      <ErrorMessage message="Heads up" type="warning" />
    );
    expect(screen.getByRole("alert")).toHaveClass("alert-warning");

    rerender(
      <ErrorMessage message="FYI" type="info" />
    );
    expect(screen.getByRole("alert")).toHaveClass("alert-info");
  });

  it("calls onDismiss when the close button is clicked", () => {
    const onDismiss = vi.fn();
    renderWithProviders(
      <ErrorMessage message="Dismiss me" onDismiss={onDismiss} />
    );

    fireEvent.click(screen.getByRole("button", { name: /close/i }));
    expect(onDismiss).toHaveBeenCalledTimes(1);
  });

  it("omits the dismiss button when onDismiss is not provided", () => {
    renderWithProviders(<ErrorMessage message="No dismiss" />);
    expect(screen.queryByRole("button", { name: /close/i })).not.toBeInTheDocument();
  });
});
