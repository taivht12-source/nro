package combine;
import consts.ConstNpc;
import item.Item;
import player.Player;
import player.Service.InventoryService;
import services.ItemService;
import services.Service;
import utils.Util;

public class NangCapBongTai {

    private static final int GOLD_BONG_TAI = 200_000_000;
    private static final int GEM_BONG_TAI = 1_000;
    private static final int GEM_NANG_BT = 1_000;
    private static final int RATIO_BONG_TAI = 50;

    public static void showInfoCombine(Player player) {
        if (player.combineNew.itemsCombine.size() == 2) {
            Item bongTai = null;
            Item manhVo = null;
            for (Item item : player.combineNew.itemsCombine) {
                if (item.template.id == 454) {
                    bongTai = item;
                } else if (item.template.id == 933) {
                    manhVo = item;
                }
            }
            if (bongTai != null && manhVo != null) {

                player.combineNew.goldCombine = GOLD_BONG_TAI;
                player.combineNew.gemCombine = GEM_BONG_TAI;
                player.combineNew.ratioCombine = RATIO_BONG_TAI;

                String npcSay = "|2|Bông tai Porata [+2]" + "\n\n";
                npcSay += "|2|Tỉ lệ thành công: " + player.combineNew.ratioCombine + "%" + "\n";
                if (InventoryService.gI().getParam(player, 31, 933) < 9999) {
                    npcSay += "|7|Cần 9999 " + manhVo.template.name + "\n";
                    npcSay += "|2|Cần: " + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n";
                    npcSay += "|2|Cần: " + player.combineNew.gemCombine + " ngọc\n";
                    npcSay += "|7|Thất bại -99 " + manhVo.template.name + "\n";
                    npcSay += "Còn thiếu " + (9999 - InventoryService.gI().getParam(player, 31, 933)) + " " + manhVo.template.name;
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, npcSay, "Đóng");
                } else if (player.inventory.getGem() >= player.combineNew.gemCombine && player.inventory.gold >= player.combineNew.goldCombine) {
                    npcSay += "|2|Cần 9999 " + manhVo.template.name + "\n";
                    npcSay += "|2|Cần: " + Util.numberToMoney(player.combineNew.gemCombine) + " ngọc\n";
                    npcSay += "|2|Cần: " + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n";
                    npcSay += "|7|Thất bại -99 " + manhVo.template.name + "\n";
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, npcSay,
                            "Nâng cấp\n"
                            + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n"
                            + Util.numberToMoney(player.combineNew.gemCombine) + " ngọc\n", "Từ chối");
                } else if (player.inventory.gem < player.combineNew.gemCombine) {
                    npcSay += "|2|Cần 9999 " + manhVo.template.name + "\n";
                    npcSay += "|7|Cần: " + player.combineNew.gemCombine + " ngọc xanh\n";
                    npcSay += "|2|Cần: " + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n";
                    npcSay += "|7|Thất bại -99 " + manhVo.template.name + "\n";
                    npcSay += "Còn thiếu\n" + (player.combineNew.gemCombine - player.inventory.gem) + " ngọc xanh";
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, npcSay, "Đóng");
                } else if (player.inventory.gold < player.combineNew.goldCombine) {
                    npcSay += "|2|Cần 9999 " + manhVo.template.name + "\n";
                    npcSay += "|2|Cần: " + player.combineNew.gemCombine + " ngọc\n";
                    npcSay += "|7|Cần: " + Util.numberToMoney(player.combineNew.goldCombine) + " vàng\n";
                    npcSay += "|7|Thất bại -99 " + manhVo.template.name + "\n";
                    npcSay += "Còn thiếu " + Util.powerToString(player.combineNew.goldCombine - player.inventory.gold) + " vàng";
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, npcSay, "Đóng");
                }
            } else {
                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                        "Cần 1 Bông tai Porata cấp 1 và Mảnh vỡ bông tai", "Đóng");
            }
        } else {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                    "Cần 1 Bông tai Porata cấp 1 và Mảnh vỡ bông tai", "Đóng");
        }
    }

    public static void nangCapBongTai(Player player) {
        if (player.combineNew.itemsCombine.size() == 2) {
            int gold = player.combineNew.goldCombine;
            if (player.inventory.gold < gold) {
                Service.gI().sendThongBao(player, "Không đủ vàng để thực hiện");
                return;
            }
            int gem = player.combineNew.gemCombine;
            if (player.inventory.gem < gem) {
                Service.gI().sendThongBao(player, "Không đủ ngọc để thực hiện");
                return;
            }
            Item bongTai = null;
            Item manhVo = null;
            for (Item item : player.combineNew.itemsCombine) {
                if (item.template.id == 454) {
                    bongTai = item;
                } else if (item.template.id == 933) {
                    manhVo = item;
                }
            }
            if (bongTai != null && manhVo != null) {
                Item findItemBag = InventoryService.gI().findItemBag(player, 921); //Khóa btc2
                if (findItemBag != null) {
                    Service.gI().sendThongBao(player, "Ngươi đã có bông tai Porata cấp 2 trong hàng trang rồi, không thể nâng cấp nữa.");
                    return;
                }
                player.inventory.gold -= gold;
                player.inventory.gem -= gem;
                if (Util.isTrue(player.combineNew.ratioCombine, 100)) {
                    bongTai.template = ItemService.gI().getTemplate(921);
                    bongTai.itemOptions.clear();
                    bongTai.itemOptions.add(new Item.ItemOption(72, 2));
                    CombineService.gI().sendEffectSuccessCombine(player);
                    InventoryService.gI().subParamItemsBag(player, 933, 31, 9999);
                } else {
                    CombineService.gI().sendEffectFailCombine(player);
                    InventoryService.gI().subParamItemsBag(player, 933, 31, 99);
                }
                InventoryService.gI().sendItemBags(player);
                Service.gI().sendMoney(player);
                CombineService.gI().reOpenItemCombine(player);
            }
        }
    }

}
