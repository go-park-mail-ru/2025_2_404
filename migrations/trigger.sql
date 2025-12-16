-- Триггеры для обновления updated_at

CREATE TRIGGER tr_client_updated_at
    BEFORE UPDATE ON client
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_client_wallet_updated_at
    BEFORE UPDATE ON client_wallet
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_wallet_top_up_updated_at
    BEFORE UPDATE ON wallet_top_up
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_ad_updated_at
    BEFORE UPDATE ON ad
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_ad_detail_updated_at
    BEFORE UPDATE ON ad_detail
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_statistic_updated_at
    BEFORE UPDATE ON statistic
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_slots_updated_at
    BEFORE UPDATE ON slots
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_create_wallet_after_client_insert
    AFTER INSERT ON client
    FOR EACH ROW
    EXECUTE FUNCTION create_wallet_for_new_client();

CREATE TRIGGER trigger_create_statistic
    AFTER INSERT ON ad_detail
    FOR EACH ROW
    EXECUTE FUNCTION create_statistic_for_ad_detail();